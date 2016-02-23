// Copyright 2015 crane authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tsuru/tsuru/cmd"
	"github.com/tsuru/tsuru/cmd/cmdtest"
	"gopkg.in/check.v1"
)

func (s *S) TestServiceCreateInfo(c *check.C) {
	desc := "Creates a service based on a passed manifest. The manifest format should be a yaml and follow the standard described in the documentation (should link to it here)"
	cmd := ServiceCreate{}
	i := cmd.Info()
	c.Assert(i.Name, check.Equals, "create")
	c.Assert(i.Usage, check.Equals, "create path/to/manifest [- for stdin]")
	c.Assert(i.Desc, check.Equals, desc)
	c.Assert(i.MinArgs, check.Equals, 1)
}

func (s *S) TestServiceCreateRun(c *check.C) {
	var stdout, stderr bytes.Buffer
	args := []string{"testdata/manifest.yml"}
	context := cmd.Context{
		Args:   args,
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := cmdtest.Transport{Message: "success", Status: http.StatusOK}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	err := (&ServiceCreate{}).Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, "Service successfully created\n")
}

func (s *S) TestServiceRemoveRun(c *check.C) {
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	stdin := bytes.NewBufferString("y\n")
	context := cmd.Context{
		Args:   []string{"my-service"},
		Stdout: &stdout,
		Stderr: &stderr,
		Stdin:  stdin,
	}
	trans := cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{
			Message: "",
			Status:  http.StatusNoContent,
		},
		CondFunc: func(req *http.Request) bool {
			called = true
			return req.Method == "DELETE" && strings.HasSuffix(req.URL.Path, "/services/my-service")
		},
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	err := (&ServiceRemove{}).Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(called, check.Equals, true)
	expected := `Are you sure you want to remove the service "my-service"? (y/n) Service successfully removed.`
	c.Assert(stdout.String(), check.Equals, expected+"\n")
}

func (s *S) TestServiceRemoveRunWithRequestFailure(c *check.C) {
	var stdout, stderr bytes.Buffer
	context := cmd.Context{
		Args:   []string{"my-service"},
		Stdout: &stdout,
		Stderr: &stderr,
		Stdin:  bytes.NewBufferString("y\n"),
	}
	trans := cmdtest.Transport{
		Message: "This service cannot be removed because it has instances.\nPlease remove these instances before removing the service.",
		Status:  http.StatusForbidden,
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	err := (&ServiceRemove{}).Run(&context, client)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, trans.Message)
}

func (s *S) TestServiceRemoveIsACommand(c *check.C) {
	var _ cmd.Command = &ServiceRemove{}
}

func (s *S) TestServiceRemoveInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "remove",
		Usage:   "remove <servicename>",
		Desc:    "removes a service from catalog",
		MinArgs: 1,
	}
	c.Assert((&ServiceRemove{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestServiceListInfo(c *check.C) {
	cmd := ServiceList{}
	i := cmd.Info()
	c.Assert(i.Name, check.Equals, "list")
	c.Assert(i.Usage, check.Equals, "list")
	c.Assert(i.Desc, check.Equals, "list services that belongs to user's team and it's service instances.")
}

func (s *S) TestServiceListRun(c *check.C) {
	var stdout, stderr bytes.Buffer
	response := `[{"service": "mysql", "instances": ["my_db"]}]`
	expected := `+----------+-----------+
| Services | Instances |
+----------+-----------+
| mysql    | my_db     |
+----------+-----------+
`
	trans := cmdtest.Transport{Message: response, Status: http.StatusOK}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	context := cmd.Context{
		Args:   []string{},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err := (&ServiceList{}).Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, expected)
}

func (s *S) TestServiceListRunWithNoServicesReturned(c *check.C) {
	var stdout, stderr bytes.Buffer
	response := `[]`
	expected := ""
	trans := cmdtest.Transport{Message: response, Status: http.StatusOK}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	context := cmd.Context{
		Args:   []string{},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err := (&ServiceList{}).Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, expected)
}

func (s *S) TestServiceUpdate(c *check.C) {
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	trans := cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{
			Message: "",
			Status:  http.StatusNoContent,
		},
		CondFunc: func(req *http.Request) bool {
			called = true
			return req.Method == "PUT" && strings.HasSuffix(req.URL.Path, "/services")
		},
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	context := cmd.Context{
		Args:   []string{"testdata/manifest.yml"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err := (&ServiceUpdate{}).Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(called, check.Equals, true)
	c.Assert(stdout.String(), check.Equals, "Service successfully updated.\n")
}

func (s *S) TestServiceUpdateIsACommand(c *check.C) {
	var _ cmd.Command = &ServiceUpdate{}
}

func (s *S) TestServiceUpdateInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "update",
		Usage:   "update <path/to/manifest>",
		Desc:    "Update service data, extracting it from the given manifest file.",
		MinArgs: 1,
	}
	c.Assert((&ServiceUpdate{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestServiceDocAdd(c *check.C) {
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	trans := cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: "", Status: http.StatusNoContent},
		CondFunc: func(req *http.Request) bool {
			called = true
			return req.Method == "PUT" && strings.HasSuffix(req.URL.Path, "/services/serv/doc")
		},
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	context := cmd.Context{
		Args:   []string{"serv", "testdata/doc.md"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err := (&ServiceDocAdd{}).Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(called, check.Equals, true)
	c.Assert(stdout.String(), check.Equals, "Documentation for 'serv' successfully updated.\n")
}

func (s *S) TestServiceDocAddInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "doc-add",
		Usage:   "service doc-add <service> <path/to/docfile>",
		Desc:    "Update service documentation, extracting it from the given file.",
		MinArgs: 2,
	}
	c.Assert((&ServiceDocAdd{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestServiceDocGet(c *check.C) {
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	trans := cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: "some doc", Status: http.StatusNoContent},
		CondFunc: func(req *http.Request) bool {
			called = true
			return req.Method == "GET" && strings.HasSuffix(req.URL.Path, "/services/serv/doc")
		},
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	context := cmd.Context{
		Args:   []string{"serv"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err := (&ServiceDocGet{}).Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(called, check.Equals, true)
	c.Assert(context.Stdout.(*bytes.Buffer).String(), check.Equals, "some doc")
}

func (s *S) TestServiceDocGetInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "doc-get",
		Usage:   "service doc-get <service>",
		Desc:    "Shows service documentation.",
		MinArgs: 1,
	}
	c.Assert((&ServiceDocGet{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestServiceTemplateInfo(c *check.C) {
	got := (&ServiceTemplate{}).Info()
	usg := `template
e.g.: $ crane template`
	expected := &cmd.Info{
		Name:  "template",
		Usage: usg,
		Desc:  "Generates a manifest template file and places it in current directory",
	}
	c.Assert(got, check.DeepEquals, expected)
}

func (s *S) TestServiceTemplateRun(c *check.C) {
	var stdout, stderr bytes.Buffer
	trans := cmdtest.Transport{Message: "", Status: http.StatusOK}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	ctx := cmd.Context{
		Args:   []string{},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err := (&ServiceTemplate{}).Run(&ctx, client)
	defer os.Remove("./manifest.yaml")
	c.Assert(err, check.IsNil)
	expected := "Generated file \"manifest.yaml\" in current directory\n"
	c.Assert(stdout.String(), check.Equals, expected)
	f, err := os.Open("./manifest.yaml")
	c.Assert(err, check.IsNil)
	fc, err := ioutil.ReadAll(f)
	manifest := `id: servicename
username: username_to_auth
password: .{16}
team: team_responsible_to_provide_service
endpoint:
  production: production-endpoint.com`
	c.Assert(string(fc), check.Matches, manifest)
}
