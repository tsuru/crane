// Copyright 2015 crane authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/tsuru/tsuru/cmd"
)

const (
	version = "1.0.0"
	header  = "Supported-Crane"
)

func buildManager(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version, header, nil)
	m.RegisterRemoved("create", "You should use `tsuru service-create` instead.")
	m.RegisterRemoved("remove", "You should use `tsuru service-destroy` instead.")
	m.RegisterRemoved("list", "You should use `tsuru service-list` instead.")
	m.RegisterRemoved("update", "You should use `tsuru service-update` instead.")
	m.RegisterRemoved("doc-get", "You should use `tsuru service-doc-get` instead.")
	m.RegisterRemoved("doc-add", "You should use `tsuru service-doc-add` instead.")
	m.RegisterRemoved("template", "You should use `tsuru service-template` instead.")
	return m
}

func main() {
	name := cmd.ExtractProgramName(os.Args[0])
	manager := buildManager(name)
	args := os.Args[1:]
	manager.Run(args)
}
