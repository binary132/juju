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
	serviceTag names.serviceTag
	fullSchema bool
	out        cmd.Output
}

const definedDoc = `
Show the actions available to run on the target service, with a short
description.  To show the schema for the actions, use --schema.

For more information, see also the 'do' subcommand, which executes actions.
`

// Set up the YAML output.
func (c *DefinedCommand) SetFlags(f *gnuflag.FlagSet) {
	// TODO(binary132) add json output?
	c.out.AddFlags(f, "yaml", map[string]cmd.Formatter{
		"yaml": cmd.FormatYaml,
	})
	f.BoolVar(&c.fullSchema, "schema", false, "display the full action schema")
}

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
		c.serviceTag = names.NewserviceTag(svcName)
		return nil
	default:
		return cmd.CheckEmpty(args[1:])
	}
}

// Run grabs the Actions spec from the api.  It then sets up a sensible
// output format for the map.
func (c *DefinedCommand) Run(ctx *cmd.Context) error {
	api, err := c.NewActionAPIClient()
	if err != nil {
		return err
	}
	defer api.Close()

	actions, err := api.ServiceCharmActions(c.serviceTag)
	if err != nil {
		return err
	}
	actionSpecs := actions.ActionSpecs
	numActionSpecs := len(actions.ActionSpecs)
	if numActionSpecs == 0 {
		return c.out.Write(ctx, "No actions defined for %s", c.serviceTag)
	}

	if !c.fullSchema {
		tabbedResults := [][]string{}
		for name, spec := range actionSpecs.ActionSpecs {
			tabbedResults = append(tabbedResults, []string{name, spec.Description})
		}
		output, err := writeTabbedString(tabbedResults)
		if err != nil {
			return errors.Wrap(err, errors.New("action formatting failed"))
		}
		return c.out.Write(ctx, output)
	}
	return c.out.Write(ctx, actionSpecs)
}
