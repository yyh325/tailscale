# Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e
if [ -d /run/systemd/system ]; then
	deb-systemd-invoke stop 'tailscaled.service' >/dev/null || true
fi
