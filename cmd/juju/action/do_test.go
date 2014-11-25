// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action_test

import (
	"strings"

	"github.com/juju/juju/cmd/juju/action"
	"github.com/juju/juju/testing"
	"github.com/juju/names"
	gc "gopkg.in/check.v1"
)

type FetchSuite struct {
	BaseActionSuite
	subcommand *action.FetchCommand
}

var _ = gc.Suite(&FetchSuite{})

func (s *FetchSuite) SetUpTest(c *gc.C) {
	s.BaseActionSuite.SetUpTest(c)
}

func (s *FetchSuite) TestHelp(c *gc.C) {
	s.checkHelp(c, s.subcommand)
}

func (s *FetchSuite) TestInit(c *gc.C) {
	tests := []struct {
		should               string
		args                 []string
		expectUnit           names.UnitTag
		expectAction         string
		expectAsync          bool
		expectParamsYamlPath string
		expectOutput         string
		expectError          string
	}{{
		should:      "fail with missing args",
		args:        []string{},
		expectError: "no action UUID specified",
	}, {
		should:      "fail with no action specified",
		args:        []string{validUnitId},
		expectError: "no action specified",
	}, {
		should:      "fail with invalid unit tag",
		args:        []string{invalidUnitId},
		expectError: "invalid unit name \"poop\"",
	}, {
		should:      "fail with invalid action name",
		args:        []string{validUnitId, "BadName"},
		expectError: "oops poops",
	}, {
		should:      "fail with too many args",
		args:        []string{"1", "2", "3"},
		expectError: "too many poops",
	}, {
		should:       "init properly with no params",
		args:         []string{validUnitId, "valid-action-name"},
		expectUnit:   names.NewUnitTag(validUnitId),
		expectAction: "valid-action-name",
	}, {
		should:      "handle --async properly",
		args:        []string{"--async", validUnitId, "valid-action-name"},
		expectAsync: true,
	}, {
		should:      "handle --params properly",
		args:        []string{"--async", validUnitId, "valid-action-name"},
		expectAsync: true,
	}, {
		should: "handle both --params and --async properly",
		args: []string{"--async", "--params=somefile.yaml",
			validUnitId, "valid-action-name"},
		expectAsync:          true,
		expectParamsYamlPath: "somefile.yaml",
	}}

	for i, t := range tests {
		s.subcommand = &action.DoCommand{}
		c.Logf("test %d: it should %s: juju actions do %s", i,
			t.should, strings.Join(t.args, " "))
		err := testing.InitCommand(s.subcommand, t.args)
		if t.expectError == "" {
			c.Check(s.subcommand.UnitTag(), gc.Equals, t.expectUnit)
			c.Check(s.subcommand.ActionName(), gc.Equals, t.expectAction)
			c.Check(s.subcommand.IsAsync(), gc.Equals, t.expectAsync)
			c.Check(s.subcommand.ParamsYAMLPath(), gc.Equals, t.expectParamsYamlPath)
		} else {
			c.Check(err, gc.ErrorMatches, t.expectError)
		}
	}
}

// func (s *FetchSuite) TestRun(c *gc.C) {
// 	tests := []struct {
// 		should         string
// 		withResults    []params.ActionResult
// 		withAPIError   string
// 		expectedErr    string
// 		expectedOutput string
// 	}{{
// 		should:       "pass api error through properly",
// 		withAPIError: "api call error",
// 		expectedErr:  "api call error",
// 	}, {
// 		should:         "fail gracefully with no results",
// 		withResults:    []params.ActionResult{},
// 		expectedOutput: "No results for action \"action-service-name/0_a_0\"\n",
// 	}, {
// 		should:      "error correctly with multiple results",
// 		withResults: []params.ActionResult{{}, {}},
// 		expectedErr: "too many results for action \"action-service-name/0_a_0\"",
// 	}, {
// 		should: "pass through an error from the API server",
// 		withResults: []params.ActionResult{{
// 			Error: common.ServerError(errors.New("an apiserver error")),
// 		}},
// 		expectedErr: "an apiserver error",
// 	}, {
// 		should: "pretty-print action output",
// 		withResults: []params.ActionResult{{
// 			Status:  "complete",
// 			Message: "oh dear",
// 			Output: map[string]interface{}{
// 				"foo": map[string]interface{}{
// 					"bar": "baz",
// 				},
// 			},
// 		}},
// 		expectedOutput: "message: oh dear\n" +
// 			"results:\n" +
// 			"  foo:\n" +
// 			"    bar: baz\n" +
// 			"status: complete\n",
// 	}}
//
// 	for i, t := range tests {
// 		func() { // for the defer of restoring patch function
// 			s.subcommand = &action.FetchCommand{}
// 			c.Logf("test %d: it should %s", i, t.should)
// 			client := &fakeAPIClient{
// 				actionResults: t.withResults,
// 			}
// 			if t.withAPIError != "" {
// 				client.apiErr = errors.New(t.withAPIError)
// 			}
// 			restore := s.BaseActionSuite.patchAPIClient(client)
// 			defer restore()
// 			//args := fmt.Sprintf("%s %s", s.subcommand.Info().Name, "some-action-id")
// 			ctx, err := testing.RunCommand(c, s.subcommand, validActionId)
// 			if t.expectedErr != "" || t.withAPIError != "" {
// 				c.Check(err, gc.ErrorMatches, t.expectedErr)
// 			} else {
// 				c.Assert(err, gc.IsNil)
// 				c.Check(ctx.Stdout.(*bytes.Buffer).String(), gc.Matches, t.expectedOutput)
// 			}
// 		}()
// 	}
// }
