// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action

import (
	"fmt"
	"regexp"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/names"
	"launchpad.net/gnuflag"
)

// DoCommand enqueues an Action for running on the given unit with given
// params
type DoCommand struct {
	ActionCommandBase
	unitTag      names.UnitTag
	actionName   string
	actionParams map[string]interface{}
	ParamsYAML   cmd.FileVar
	out          cmd.Output
	undefinedActionCommand
}

const doDoc = `
Queue an Action for execution on a given unit, with a given set of params.
Displays the ID of the Action for use with 'juju kill', 'juju status', etc.
The command will wait until it receives a result unless --async is used.

Params are validated according to the charm for the unit's service.  The 
valid params can be seen using "juju action defined <service>".  Params must
be in a yaml file which is passed with the --params flag.

Examples:

$ juju do mysql/2 pause

finished

$ juju do mysql/3 backup --async
action: <UUID>

$ juju status <UUID>
result:
  status: success
  file:
    size: 873.2
    units: GB
    name: foo.sql

$ juju do mysql/3 backup --async --params parameters.yml
...
`

// SetFlags offers an option for YAML output.
func (c *DoCommand) SetFlags(f *gnuflag.FlagSet) {
	// TODO(binary132) add json output?
	c.out.AddFlags(f, "yaml", map[string]cmd.Formatter{
		"yaml": cmd.FormatYaml,
	})
	f.Var(&c.ParamsYAML, "params", "path to yaml-formatted params file")
}

func (c *DoCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "do",
		Args:    "<unit> <action name>",
		Purpose: "WIP: queue an action for execution",
		Doc:     doDoc,
	}
}

// Init gets the unit tag, and checks for other correct args.
func (c *DoCommand) Init(args []string) error {
	switch len(args) {
	case 0:
		return errors.New("no unit specified")
	case 1:
		return errors.New("no action specified")
	case 2:
		unitName := args[0]
		if !names.IsValidUnit(unitName) {
			return errors.Errorf("invalid unit name %q", unitName)
		}
		actionName := args[1]
		actionNameRule := regexp.MustCompile("^[a-z](?:[a-z-]*[a-z])?$")
		if valid := actionNameRule.MatchString(actionName); !valid {
			return fmt.Errorf("invalid action name %q", actionName)
		}
		c.unitTag = names.NewUnitTag(unitName)
		c.actionName = actionName

		return nil
	default:
		return cmd.CheckEmpty(args[1:])
	}
}

func (c *DoCommand) Run(ctx *cmd.Context) error {
	api, err := c.NewActionAPIClient()
	if err != nil {
		return err
	}
	defer api.Close()

	c.actionParams = map[string]interface{}{}

	if c.ParamsYAML.Path != "" {
		b, err := c.ParamsYAML.Read(ctx)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(b, &c.actionParams)

		conformantParams, err := conform(c.actionParams)
		if err != nil {
			return err
		}

		betterParams, ok := conformantParams.(map[string]interface{})
		if !ok {
			return errors.New("params must contain a YAML map with string keys")
		}

		c.actionParams = betterParams
	}

	actionParam := params.Actions{
		Actions: []params.Action{{
			Receiver:   c.unitTag,
			Name:       c.actionName,
			Parameters: c.actionParams,
		}},
	}

	result, err := api.Enqueue(actionParam)
	if err != nil {
		return err
	}
	if len(result.Results) != 1 {
		return errors.New("only one result must be received")
	}

	err = result.Results[0].Error
	if err != nil {
		return err
	}

	tag := result.Results[0].Action.Tag
	if !names.IsValidAction(tag.Id()) {
		return errors.Errorf("invalid action tag %q received", tag.String())
	}

	err = c.out.Write(ctx, fmt.Sprintf("Action queued with id: %#v", tag.String()))
	if err != nil {
		return err
	}

	for _ = range time.Tick(1 * time.Second) {
		completed, err := api.ListCompleted(params.Tags{Tags: []names.Tag{c.unitTag}})
		if err != nil {
			return err
		}

		if len(completed.Actions) != 1 {
			return errors.New("only one result must be received")
		}
		err = completed.Actions[0].Error
		if err != nil {
			return err
		}

		results := completed.Actions[0].Actions
		if len(results) == 0 {
			continue
		}
		if len(results) > 1 {
			return errors.New("too many action results")
		}

		err = displayActionResult(results[0], ctx, c.out)
		if err != nil {
			return err
		}
	}

	return nil
}
