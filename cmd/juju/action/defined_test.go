// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action_test

import (
	"bytes"
	"strings"

	"github.com/juju/juju/cmd/juju/action"
	"github.com/juju/juju/state"
	"github.com/juju/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	yaml "gopkg.in/yaml.v2"
)

type DefinedSuite struct {
	BaseActionSuite
	svc        *state.Service
	subcommand *action.DefinedCommand
}

var _ = gc.Suite(&DefinedSuite{})

func (s *DefinedSuite) SetUpTest(c *gc.C) {
	s.BaseActionSuite.SetUpTest(c)
	s.subcommand = &action.DefinedCommand{}
}

func (s *DefinedSuite) TestHelp(c *gc.C) {
	s.checkHelp(c, s.subcommand)
}

func (s *DefinedSuite) TestInit(c *gc.C) {
	tests := []struct {
		should      string
		args        []string
		svcTag      names.Service
		errorString string
	}{{
		should:      "fail with missing service name",
		args:        []string{},
		errorString: "no service name specified",
	}, {
		should:      "fail with invalid service name",
		args:        []string{"derp/0"},
		errorString: "invalid service name \"derp\"",
	}, {
		should: "init properly with valid service name",
		args:   []string{"mysql"},
		svcTag: names.NewServiceTag("mysql"),
	}, {
		should: "init properly with valid service name and --schema",
		args:   []string{"mysql"},
		svcTag: names.NewServiceTag("mysql"),
	}}

	for i, test := range tests {
		c.Logf("test %d should %s: juju actions defined %s", i,
			t.should, strings.Join(args, " "))
		err := testing.InitCommand(s.subcommand, t.args)
		if test.ErrorString == "" {
			c.Check(definedCmd.ServiceTag, gc.Equals, t.svcTag)
		} else {
			c.Check(err, gc.ErrorMatches, t.errorString)
		}
		// wip
	}
}

func (s *DefinedSuite) TestRun(c *gc.C) {
	tests := []struct {
		args            []string
		expectedResults map[string]interface{}
		expectedErr     string
	}{{
		args:            []string{},
		expectedErr:     "error: no service name specified\n",
		expectedResults: map[string]interface{}{},
	}, {
		args: []string{"dummy"},
		expectedResults: map[string]interface{}{
			"snapshot": map[string]interface{}{
				"description": "Take a snapshot of the database.",
				"params": map[string]interface{}{
					"type": "object",
					"outfile": map[string]interface{}{
						"default":     "foo.bz2",
						"type":        "string",
						"description": "The file to write out to.",
					},
				},
			},
		},
	}, {
		args:            []string{"dne"},
		expectedErr:     "error: service \"dne\" not found\n",
		expectedResults: map[string]interface{}{},
	}, {
		args:            []string{"two", "things"},
		expectedErr:     "error: unrecognized args: [\"things\"]\n",
		expectedResults: map[string]interface{}{},
	}, {
		args:            []string{"dummy", "things"},
		expectedErr:     "error: unrecognized args: [\"things\"]\n",
		expectedResults: map[string]interface{}{},
	}}

	ch := s.AddTestingCharm(c, "dummy")
	svc := s.AddTestingService(c, "dummy", ch)
	s.svc = svc

	for i, t := range tests {
		c.Logf("test %d: %#v", i, t.args)
		args := append(s.subcommand.Name, t.args)
		ctx, err := testing.RunCommand(c, s.command, args)
		c.Assert(err, gc.IsNil)
		buf, err := yaml.Marshal(t.expectedResults)
		s.checkStd(c, ctx, t.buf, t.expectedErr)
		// wip
		expected := make(map[string]interface{})
		err = yaml.Unmarshal(buf, &expected)
		c.Assert(err, gc.IsNil)
		actual := make(map[string]interface{})
		err = yaml.Unmarshal(ctx.Stdout.(*bytes.Buffer).Bytes(), &actual)
		c.Assert(err, gc.IsNil)
		c.Check(actual, jc.DeepEquals, expected)
	}
}
