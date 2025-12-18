// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// The hello server contains a single tool that says hi to the user.
//
// It runs over the stdio transport.
package main

import (
	"github.com/aidanuno/qrkit/cmd"
)

func main() {
	cmd.Start()
}
