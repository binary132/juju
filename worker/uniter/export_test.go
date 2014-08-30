// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"github.com/juju/names"
	"github.com/juju/utils/proxy"
)

func SetUniterObserver(u *Uniter, observer UniterExecutionObserver) {
	u.observer = observer
}

func (u *Uniter) GetProxyValues() proxy.Settings {
	u.proxyMutex.Lock()
	defer u.proxyMutex.Unlock()
	return u.proxy
}

func (c *HookContext) ActionResultsMap() map[string]interface{} {
	if c.actionData != nil {
		return c.actionData.ResultsMap
	}
	return map[string]interface{}{}
}

func (c *HookContext) ActionFailed() bool {
	return c.actionData.ActionFailed
}

func (c *HookContext) ActionMessage() string {
	return c.actionData.ResultsMessage
}

func (c *HookContext) SetActionData(tag names.ActionTag, params map[string]interface{}) {
	c.actionData = &actionData{
		ActionTag:    tag,
		ActionParams: params,
		ResultsMap:   map[string]interface{}{},
	}
}

var MergeEnvironment = mergeEnvironment

var SearchHook = searchHook

var HookCommand = hookCommand

var LookPath = lookPath
