// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package jujuc_test

import (
	"github.com/juju/cmd"
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
		command     []string
		failMessage string
		errMsg      string
		code        int
	}{{
		summary:     "no parameters sets a default message",
		command:     []string{},
		failMessage: "action failed without reason given, check action for errors",
	}}

	for i, t := range actionFailTests {
		c.Logf("test %d: %s", i, t.summary)
		hctx := &Context{}
		com, err := jujuc.NewCommand(hctx, "action-fail")
		c.Assert(err, gc.IsNil)
		ctx := testing.Context(c)
		code := cmd.Main(com, ctx, t.command)
		c.Check(code, gc.Equals, t.code)
		_, actionErr := hctx.ActionResults()
		c.Check(bufferString(ctx.Stderr), gc.Equals, t.errMsg)
		c.Check(actionErr, gc.ErrorMatches, t.failMessage)
	}
}

func (s *ActionFailSuite) TestHelp(c *gc.C) {
	hctx := s.GetHookContext(c, -1, "")
	com, err := jujuc.NewCommand(hctx, "action-fail")
	c.Assert(err, gc.IsNil)
	ctx := testing.Context(c)
	code := cmd.Main(com, ctx, []string{"--help"})
	c.Assert(code, gc.Equals, 0)
	c.Assert(bufferString(ctx.Stdout), gc.Equals, `usage: action-fail [options] ["<failure message>"]
purpose: set action fail status with message

options:
--clear  (= false)
    clear an existing fail state

action-fail sets the action's fail state with a given error message.  Using
action-fail without a failure message, and without --clear, will set a
default failure message indicating a problem with the action.
`)
	c.Assert(bufferString(ctx.Stderr), gc.Equals, "")
}
