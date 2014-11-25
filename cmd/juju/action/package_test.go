// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action_test

import (
	"testing"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/cmd/juju/action"
	coretesting "github.com/juju/juju/testing"
	"github.com/juju/names"
	jujutesting "github.com/juju/testing"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/charm.v4"
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

	expected := "(?sm).*^usage: juju action " +
		subcmd.Info().Name + " " +
		"\\[options\\] " +
		subcmd.Info().Args + ".+"
	c.Check(coretesting.Stdout(ctx), gc.Matches, expected)

	expected = "(?sm).*^purpose: " + subcmd.Info().Purpose + "$.*"
	c.Check(coretesting.Stdout(ctx), gc.Matches, expected)

	expected = "(?sm).*^" + subcmd.Info().Doc + "$.*"
	c.Check(coretesting.Stdout(ctx), gc.Matches, expected)
}

type fakeAPIClient struct {
	action.APIClient
	actionResults      []params.ActionResult
	actionsByReceivers []params.ActionsByReceiver
	charmActions       *charm.Actions
	apiErr             error
}

func (c *fakeAPIClient) Close() error {
	return nil
}

func (c *fakeAPIClient) Enqueue(params.Actions) (params.ActionResults, error) {
	return params.ActionResults{
		Results: c.actionResults,
	}, c.apiErr
}

func (c *fakeAPIClient) ListAll(args params.Entities) (params.ActionsByReceivers, error) {
	return params.ActionsByReceivers{
		Actions: c.actionsByReceivers,
	}, c.apiErr
}

func (c *fakeAPIClient) ListPending(args params.Entities) (params.ActionsByReceivers, error) {
	return params.ActionsByReceivers{
		Actions: c.actionsByReceivers,
	}, c.apiErr
}

func (c *fakeAPIClient) ListCompleted(args params.Entities) (params.ActionsByReceivers, error) {
	return params.ActionsByReceivers{
		Actions: c.actionsByReceivers,
	}, c.apiErr
}

func (c *fakeAPIClient) Cancel(args params.Actions) (params.ActionResults, error) {
	return params.ActionResults{
		Results: c.actionResults,
	}, c.apiErr
}

func (c *fakeAPIClient) ServiceCharmActions(names.ServiceTag) (*charm.Actions, error) {
	return c.charmActions, c.apiErr
}

func (c *fakeAPIClient) Actions(args params.Entities) (params.ActionResults, error) {
	return params.ActionResults{
		Results: c.actionResults,
	}, c.apiErr
}
