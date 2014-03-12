// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxc

import (
	"fmt"

	"launchpad.net/juju-core/container"
	"launchpad.net/juju-core/utils"
)

var requiredPackages = []string{
	"lxc",
}

type containerInitialiser struct {
	targetRelease string
}

// containerInitialiser implements container.Initialiser.
var _ container.Initialiser = (*containerInitialiser)(nil)

// NewContainerInitialiser returns an instance used to perform the steps
// required to allow a host machine to run a LXC container.
func NewContainerInitialiser(targetRelease string) container.Initialiser {
	return &containerInitialiser{targetRelease}
}

// Initialise is specified on the container.Initialiser interface.
func (ci *containerInitialiser) Initialise() error {
	return ensureDependencies((*ci).targetRelease)
}

// targetReleasePackages returns a slice suitable for use
// with utils.AptGetInstall based on the targetRelease string.
func targetReleasePackages(targetRelease string) []string {
	packages := []string{
		fmt.Sprintf("--target-release '%s' '%s'", targetRelease, requiredPackages[0]),
	}
	return packages
}

// ensureDependencies checks the targetRelease and updates the packages
// that are sent to utils.AptGetInstall to include the --target-release
// switch. If targetRelease is an empty string, no switch is passed.
func ensureDependencies(targetRelease string) error {
	var packages []string
	if targetRelease != "" {
		packages = targetReleasePackages(targetRelease)
	} else {
		packages = requiredPackages
	}
	return utils.AptGetInstall(packages...)
}
