// Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wgengine

import (
	"errors"

	"inet.af/netaddr"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tailcfg"
	"tailscale.com/types/netmap"
	"tailscale.com/wgengine/filter"
	"tailscale.com/wgengine/monitor"
	"tailscale.com/wgengine/router"
	"tailscale.com/wgengine/tsdns"
	"tailscale.com/wgengine/wgcfg"
)

// Status is the Engine status.
//
// TODO(bradfitz): remove this, subset of ipnstate? Need to migrate users.
type Status struct {
	Peers      []ipnstate.PeerStatusLite
	LocalAddrs []string // the set of possible endpoints for the magic conn
	DERPs      int      // number of active DERP connections
}

// StatusCallback is the type of status callbacks used by
// Engine.SetStatusCallback.
//
// Exactly one of Status or error is non-nil.
type StatusCallback func(*Status, error)

// NetInfoCallback is the type used by Engine.SetNetInfoCallback.
type NetInfoCallback func(*tailcfg.NetInfo)

// NetworkMapCallback is the type used by callbacks that hook
// into network map updates.
type NetworkMapCallback func(*netmap.NetworkMap)

// someHandle is allocated so its pointer address acts as a unique
// map key handle. (It needs to have non-zero size for Go to guarantee
// the pointer is unique.)
type someHandle struct{ _ byte }

// ErrNoChanges is returned by Engine.Reconfig if no changes were made.
var ErrNoChanges = errors.New("no changes made to Engine config")

// Engine is the Tailscale WireGuard engine interface.
type Engine interface {
	// Reconfig reconfigures WireGuard and makes sure it's running.
	// This also handles setting up any kernel routes.
	//
	// This is called whenever tailcontrol (the control plane)
	// sends an updated network map.
	//
	// The returned error is ErrNoChanges if no changes were made.
	Reconfig(*wgcfg.Config, *router.Config) error

	// GetFilter returns the current packet filter, if any.
	GetFilter() *filter.Filter

	// SetFilter updates the packet filter.
	SetFilter(*filter.Filter)

	// SetDNSMap updates the DNS map.
	SetDNSMap(*tsdns.Map)

	// SetStatusCallback sets the function to call when the
	// WireGuard status changes.
	SetStatusCallback(StatusCallback)

	// GetLinkMonitor returns the link monitor.
	GetLinkMonitor() *monitor.Mon

	// RequestStatus requests a WireGuard status update right
	// away, sent to the callback registered via SetStatusCallback.
	RequestStatus()

	// Close shuts down this wireguard instance, remove any routes
	// it added, etc. To bring it up again later, you'll need a
	// new Engine.
	Close()

	// Wait waits until the Engine's Close method is called or the
	// engine aborts with an error. You don't have to call this.
	// TODO: return an error?
	Wait()

	// SetDERPMap controls which (if any) DERP servers are used.
	// If nil, DERP is disabled. It starts disabled until a DERP map
	// is configured.
	SetDERPMap(*tailcfg.DERPMap)

	// SetNetworkMap informs the engine of the latest network map
	// from the server. The network map's DERPMap field should be
	// ignored as as it might be disabled; get it from SetDERPMap
	// instead.
	// The network map should only be read from.
	SetNetworkMap(*netmap.NetworkMap)

	// AddNetworkMapCallback adds a function to a list of callbacks
	// that are called when the network map updates. It returns a
	// function that when called would remove the function from the
	// list of callbacks.
	AddNetworkMapCallback(NetworkMapCallback) (removeCallback func())

	// SetNetInfoCallback sets the function to call when a
	// new NetInfo summary is available.
	SetNetInfoCallback(NetInfoCallback)

	// DiscoPublicKey gets the public key used for path discovery
	// messages.
	DiscoPublicKey() tailcfg.DiscoKey

	// UpdateStatus populates the network state using the provided
	// status builder.
	UpdateStatus(*ipnstate.StatusBuilder)

	// Ping is a request to start a discovery ping with the peer handling
	// the given IP and then call cb with its ping latency & method.
	Ping(ip netaddr.IP, cb func(*ipnstate.PingResult))
}
