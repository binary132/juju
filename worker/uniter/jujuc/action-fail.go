// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package jujuc

import (
	"fmt"

	"github.com/juju/cmd"
	"launchpad.net/gnuflag"
)

// ActionFailCommand implements the action-fail command.
type ActionFailCommand struct {
	cmd.CommandBase
	ctx         Context
	clear       bool
	failMessage string
}

// NewActionFailCommand returns a new ActionFailCommand with the given context.
func NewActionFailCommand(ctx Context) cmd.Command {
	return &ActionFailCommand{ctx: ctx}
}

// Info returns the content for --help.
func (c *ActionFailCommand) Info() *cmd.Info {
	doc := `
action-fail sets the action's fail state with a given error message.  Using
action-fail without a failure message, and without --clear, will set a
default failure message indicating a problem with the action.
`
	return &cmd.Info{
		Name:    "action-fail",
		Args:    "[\"<failure message>\"]",
		Purpose: "set action fail status with message",
		Doc:     doc,
	}
}

// SetFlags handles option flags.  --clear will reset the fail state.
func (c *ActionFailCommand) SetFlags(f *gnuflag.FlagSet) {
	f.BoolVar(&c.clear, "clear", false, "clear an existing fail state")
}

// Init sets the fail message and checks for malformed invocations.
func (c *ActionFailCommand) Init(args []string) error {
	if c.clear {
		return cmd.CheckEmpty(args)
	}

	switch len(args) {
	case 0:
		c.failMessage = "action failed without reason given, check action for errors"
	case 1:
		c.failMessage = args[0]
	default:
		return cmd.CheckEmpty(args[1:])
	}

	return nil
}

// Run sets the Action's fail state, or clears it if --clear was passed.
func (c *ActionFailCommand) Run(ctx *cmd.Context) error {
	err := fmt.Errorf(c.failMessage)

	if c.clear {
		err = nil
	}

	c.ctx.ActionSetFailed(err)
	return nil
}
