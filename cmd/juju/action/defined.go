// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action

import (
	"github.com/juju/cmd"
	errors "github.com/juju/errors"
	"github.com/juju/names"
	"launchpad.net/gnuflag"
)

// DefinedCommand lists actions defined by the charm of a given service.
type DefinedCommand struct {
	ActionCommandBase
	ServiceTag names.ServiceTag
	out        cmd.Output
}

const definedDoc = `
Show the actions available to run on the target service, with a short
description.  To show the schema for the actions, use --schema.

For more information, see also the 'do' subcommand, which executes actions.
`

func (c *DefinedCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "defined",
		Args:    "[--schema] <service name>",
		Purpose: "WIP: show actions defined for a service",
		Doc:     definedDoc,
	}
}

// Init validates the service name and any other options.
func (c *DefinedCommand) Init(args []string) error {
	switch len(args) {
	case 0:
		return errors.New("no service name specified")
	case 1:
		svcName := args[0]
		if !names.IsValidService(svcName) {
			return errors.Errorf("invalid service name %q", svcName)
		}
		c.ServiceTag = names.NewServiceTag(svcName)
		return nil
	default:
		return cmd.CheckEmpty(args[1:])
	}
}

// Set up the YAML output.
func (c *DefinedCommand) SetFlags(f *gnuflag.FlagSet) {
	// TODO(binary132) add json output?
	c.out.AddFlags(f, "yaml", map[string]cmd.Formatter{
		"yaml": cmd.FormatYaml,
	})
}

// Run grabs the Actions spec from the api.  It then sets up a sensible
// output format for the map.
func (c *DefinedCommand) Run(ctx *cmd.Context) error {
	api, err := c.NewActionAPIClient()
	if err != nil {
		return err
	}
	defer api.Close()

	actions, err := api.ServiceCharmActions(c.ServiceTag)
	if err != nil {
		return err
	}
	actionSpecs := actions.ActionSpecs
	return c.out.Write(ctx, actionSpecs)
}
