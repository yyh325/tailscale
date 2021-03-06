// Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interfaces

import (
	"testing"
)

func TestGetState(t *testing.T) {
	st, err := GetState()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got: %#v", st)
	t.Logf("As string: %s", st)

	st2, err := GetState()
	if err != nil {
		t.Fatal(err)
	}

	if !st.Equal(st2) {
		// let's assume nobody was changing the system network interfaces between
		// the two GetState calls.
		t.Fatal("two States back-to-back were not equal")
	}

	st.RemoveTailscaleInterfaces()
	t.Logf("As string without Tailscale:\n\t%s", st)
}

func TestLikelyHomeRouterIP(t *testing.T) {
	gw, my, ok := LikelyHomeRouterIP()
	if !ok {
		t.Logf("no result")
		return
	}
	t.Logf("myIP = %v; gw = %v", my, gw)
}
