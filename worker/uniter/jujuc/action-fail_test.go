// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package jujuc_test

import (
	"github.com/juju/cmd"
	jc "github.com/juju/testing/checkers"
	gc "launchpad.net/gocheck"

	"github.com/juju/juju/testing"
	"github.com/juju/juju/worker/uniter/jujuc"
)

type ActionFailSuite struct {
	ContextSuite
}

var _ = gc.Suite(&ActionFailSuite{})

func (s *ActionFailSuite) TestActionFail(c *gc.C) {
	var actionFailTests = []struct {
		summary     string
		commands    [][]string
		failMessage string
		errMsg      string
		code        int
	}{{
		summary: "bare value(s) are an Init error",
		commands: [][]string{
			[]string{"result"},
		},
		errMsg: "error: argument \"result\" must be of the form key...=value\n",
		code:   2,
	}, {
		summary: "a response of one key to one value",
		commands: [][]string{
			[]string{"outfile=foo.bz2"},
		},
		results: map[string]interface{}{
			"outfile": "foo.bz2",
		},
	}, {
		summary: "two keys, two values",
		commands: [][]string{
			[]string{"outfile=foo.bz2", "size=10G"},
		},
		results: map[string]interface{}{
			"outfile": "foo.bz2",
			"size":    "10G",
		},
	}, {
		summary: "multiple = are ok",
		commands: [][]string{
			[]string{"outfile=foo=bz2"},
		},
		results: map[string]interface{}{
			"outfile": "foo=bz2",
		},
	}, {
		summary: "several interleaved values",
		commands: [][]string{
			[]string{"outfile.name=foo.bz2",
				"outfile.kind.util=bzip2",
				"outfile.kind.ratio=high"},
		},
		results: map[string]interface{}{
			"outfile": map[string]interface{}{
				"name": "foo.bz2",
				"kind": map[string]interface{}{
					"util":  "bzip2",
					"ratio": "high",
				},
			},
		},
	}, {
		summary: "conflicting simple values in one command result in overwrite",
		commands: [][]string{
			[]string{"util=bzip2", "util=5"},
		},
		results: map[string]interface{}{
			"util": "5",
		},
	}, {
		summary: "conflicting simple values in two commands results in overwrite",
		commands: [][]string{
			[]string{"util=bzip2"},
			[]string{"util=5"},
		},
		results: map[string]interface{}{
			"util": "5",
		},
	}, {
		summary: "conflicted map spec: {map1:{key:val}} vs {map1:val2}",
		commands: [][]string{
			[]string{"map1.key=val", "map1=val"},
		},
		results: map[string]interface{}{
			"map1": "val",
		},
	}, {
		summary: "two-invocation conflicted map spec: {map1:{key:val}} vs {map1:val2}",
		commands: [][]string{
			[]string{"map1.key=val"},
			[]string{"map1=val"},
		},
		results: map[string]interface{}{
			"map1": "val",
		},
	}, {
		summary: "conflicted map spec: {map1:val2} vs {map1:{key:val}}",
		commands: [][]string{
			[]string{"map1=val", "map1.key=val"},
		},
		results: map[string]interface{}{
			"map1": map[string]interface{}{
				"key": "val",
			},
		},
	}, {
		summary: "two-invocation conflicted map spec: {map1:val2} vs {map1:{key:val}}",
		commands: [][]string{
			[]string{"map1=val"},
			[]string{"map1.key=val"},
		},
		results: map[string]interface{}{
			"map1": map[string]interface{}{
				"key": "val",
			},
		},
	}}

	for i, t := range actionFailTests {
		c.Logf("test %d: %s", i, t.summary)
		hctx := s.GetHookContext(c, -1, "")
		com, err := jujuc.NewCommand(hctx, "action-fail")
		c.Assert(err, gc.IsNil)
		ctx := testing.Context(c)
		for j, command := range t.commands {
			c.Logf("  command %d: %#v", j, command)
			code := cmd.Main(com, ctx, command)
			c.Check(code, gc.Equals, t.code)
			results := hctx.ActionResults()
			c.Check(bufferString(ctx.Stderr), gc.Equals, t.errMsg)
			c.Check(results.errMsg, gc.Equals, t.failMessage)
			if j == len(t.commands)-1 {
				results := hctx.ActionResults()
				c.Check(results, jc.DeepEquals, t.results)
			}
		}
	}
}

func (s *ActionFailSuite) TestHelp(c *gc.C) {
	hctx := s.GetHookContext(c, -1, "")
	com, err := jujuc.NewCommand(hctx, "action-fail")
	c.Assert(err, gc.IsNil)
	ctx := testing.Context(c)
	code := cmd.Main(com, ctx, []string{"--help"})
	c.Assert(code, gc.Equals, 0)
	c.Assert(bufferString(ctx.Stdout), gc.Equals, `usage: action-fail <key>=<value> [<key>.<key>....=<value> ...]
purpose: set action results

action-fail adds the given values to the results map of the Action.  This map
is returned to the user after the completion of the Action.

Example usage:
 action-fail outfile.size=10G
 action-fail foo.bar.baz=2 foo.bar.zab=3
 action-fail foo.bar.baz=4

 will yield:

 outfile:
   size: "10G"
 foo:
   bar:
     baz: "4"
     zab: "3"
`)
	c.Assert(bufferString(ctx.Stderr), gc.Equals, "")
}
