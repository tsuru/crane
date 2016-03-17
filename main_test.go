// Copyright 2015 crane authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/tsuru/tsuru/cmd"
	"gopkg.in/check.v1"
)

func (s *S) TestCommandsFromBaseManagerAreRegistered(c *check.C) {
	baseManager := cmd.BuildBaseManager("tsuru", version, header, nil)
	manager := buildManager("tsuru")
	for name, instance := range baseManager.Commands {
		command, ok := manager.Commands[name]
		c.Assert(ok, check.Equals, true)
		c.Assert(command, check.FitsTypeOf, instance)
	}
}

func (s *S) TestCreateIsRegistered(c *check.C) {
	manager := buildManager("tsuru")
	target, ok := manager.Commands["create"]
	c.Assert(ok, check.Equals, true)
	c.Assert(target, check.FitsTypeOf, &serviceCreate{})
}

func (s *S) TestRemoveIsRegistered(c *check.C) {
	manager := buildManager("tsuru")
	remove, ok := manager.Commands["remove"]
	c.Assert(ok, check.Equals, true)
	c.Assert(remove, check.FitsTypeOf, &serviceRemove{})
}

func (s *S) TestListIsRegistered(c *check.C) {
	manager := buildManager("tsuru")
	remove, ok := manager.Commands["list"]
	c.Assert(ok, check.Equals, true)
	c.Assert(remove, check.FitsTypeOf, &serviceList{})
}

func (s *S) TestUpdateIsRegistered(c *check.C) {
	manager := buildManager("tsuru")
	update, ok := manager.Commands["update"]
	c.Assert(ok, check.Equals, true)
	c.Assert(update, check.FitsTypeOf, &serviceUpdate{})
}

func (s *S) TestDocGetIsRegistered(c *check.C) {
	manager := buildManager("tsuru")
	update, ok := manager.Commands["doc-get"]
	c.Assert(ok, check.Equals, true)
	c.Assert(update, check.FitsTypeOf, &serviceDocGet{})
}

func (s *S) TestDocAddIsRegistered(c *check.C) {
	manager := buildManager("tsuru")
	update, ok := manager.Commands["doc-add"]
	c.Assert(ok, check.Equals, true)
	c.Assert(update, check.FitsTypeOf, &serviceDocAdd{})
}

func (s *S) TestTemplateIsRegistered(c *check.C) {
	manager := buildManager("tsuru")
	update, ok := manager.Commands["template"]
	c.Assert(ok, check.Equals, true)
	c.Assert(update, check.FitsTypeOf, &serviceTemplate{})
}
