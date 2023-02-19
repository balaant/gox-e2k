package main

import (
	"fmt"
	"log"
	"strings"

	version "github.com/hashicorp/go-version"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string

	// Default, if true, will be included as a default build target
	// if no OS/arch is specified. We try to only set as a default popular
	// targets or targets that are generally useful. For example, Android
	// is not a default because it is quite rare that you're cross-compiling
	// something to Android AND something like Linux.
	Default bool
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.Arch)
}

// addDrop appends all of the "add" entries and drops the "drop" entries, ignoring
// the "Default" parameter.
func addDrop(base []Platform, add []Platform, drop []Platform) []Platform {
	newPlatforms := make([]Platform, len(base)+len(add))
	copy(newPlatforms, base)
	copy(newPlatforms[len(base):], add)

	// slow, but we only do this during initialization at most once per version
	for _, platform := range drop {
		found := -1
		for i := range newPlatforms {
			if newPlatforms[i].Arch == platform.Arch && newPlatforms[i].OS == platform.OS {
				found = i
				break
			}
		}
		if found < 0 {
			panic(fmt.Sprintf("Expected to remove %+v but not found in list %+v", platform, newPlatforms))
		}
		if found == len(newPlatforms)-1 {
			newPlatforms = newPlatforms[:found]
		} else if found == 0 {
			newPlatforms = newPlatforms[found:]
		} else {
			newPlatforms = append(newPlatforms[:found], newPlatforms[found+1:]...)
		}
	}
	return newPlatforms
}

var (
	Platforms_1_0 = []Platform{

		{"linux", "e2k", true},
	}

	// no new platforms in 1.18
	Platforms_1_18 = Platforms_1_0

	PlatformsLatest = Platforms_1_18
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	// Use latest if we get an unexpected version string
	if !strings.HasPrefix(v, "go") {
		return PlatformsLatest
	}
	// go-version only cares about version numbers
	v = v[2:]

	current, err := version.NewVersion(v)
	if err != nil {
		log.Printf("Unable to parse current go version: %s\n%s", v, err.Error())

		// Default to latest
		return PlatformsLatest
	}

	var platforms = []struct {
		constraint string
		plat       []Platform
	}{

		{">= 1.17, < 1.18", Platforms_1_18},
		{">= 1.18, < 1.19", Platforms_1_18},
	}

	for _, p := range platforms {
		constraints, err := version.NewConstraint(p.constraint)
		if err != nil {
			panic(err)
		}
		if constraints.Check(current) {
			return p.plat
		}
	}

	// Assume latest
	return PlatformsLatest
}
