// Copyright (c) 2018-present mlogclub.
// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"bbs-go/internal/server"
	_ "bbs-go/internal/services/eventhandler"
)

func main() {
	server.Init()
	server.NewServer()
}
