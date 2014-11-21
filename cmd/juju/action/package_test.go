// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action_test

import (
	"testing"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/cmd/juju/action"
	coretesting "github.com/juju/juju/testing"
	jujutesting "github.com/juju/testing"
	gc "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type BaseActionSuite struct {
	command *action.ActionCommand
}

func (s *BaseActionSuite) SetUpTest(c *gc.C) {
	s.command = action.NewActionCommand().(*action.ActionCommand)
}

func (s *BaseActionSuite) patchAPIClient(client *fakeAPIClient) func() {
	return jujutesting.PatchValue(action.NewActionAPIClient,
		func(c *action.ActionCommandBase) (action.APIClient, error) {
			return client, nil
		},
	)
}

func (s *BaseActionSuite) checkHelp(c *gc.C, subcmd envcmd.EnvironCommand) {
	ctx, err := coretesting.RunCommand(c, s.command, subcmd.Info().Name, "--help")
	c.Assert(err, gc.IsNil)

	expected := "(?s).*^usage: juju action <command> .+"
	c.Check(coretesting.Stdout(ctx), gc.Matches, expected)
	expected = "(?sm).*^purpose: " + s.command.Purpose + "$.*"
	c.Check(coretesting.Stdout(ctx), gc.Matches, expected)
	expected = "(?sm).*^" + s.command.Doc + "$.*"
	c.Check(coretesting.Stdout(ctx), gc.Matches, expected)
}

type fakeAPIClient struct {
	action.APIClient
	actionResults []params.ActionResult
	apiErr        error
}

func (c *fakeAPIClient) Close() error {
	return nil
}

func (c *fakeAPIClient) Actions(args params.ActionUUIDs) (params.ActionResults, error) {
	return params.ActionResults{
		Results: c.actionResults,
	}, c.apiErr
}
