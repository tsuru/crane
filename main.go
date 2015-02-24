// Copyright 2015 crane authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/tsuru/tsuru/cmd"
)

const (
	version = "0.7.0"
	header  = "Supported-Crane"
)

func buildManager(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version, header, nil)
	m.Register(&ServiceCreate{})
	m.Register(&ServiceRemove{})
	m.Register(&ServiceList{})
	m.Register(&ServiceUpdate{})
	m.Register(&ServiceDocGet{})
	m.Register(&ServiceDocAdd{})
	m.Register(&ServiceTemplate{})
	return m
}

func main() {
	name := cmd.ExtractProgramName(os.Args[0])
	manager := buildManager(name)
	args := os.Args[1:]
	manager.Run(args)
}
