#!/usr/bin/env sh
# Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Runs `go build` with flags configured for binary distribution. All
# it does differently from `go build` is burn git commit and version
# information into the binaries, so that we can track down user
# issues.
#
# If you're packaging Tailscale for a distro, please consider using
# this script, or executing equivalent commands in your
# distro-specific build system.

set -eu

eval $(./version/version.sh)

exec go build -tags xversion -ldflags "-X tailscale.com/version.Long=${VERSION_LONG} -X tailscale.com/version.Short=${VERSION_SHORT} -X tailscale.com/version.GitCommit=${VERSION_GIT_HASH}" "$@"
