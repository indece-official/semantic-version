package main

import "fmt"

type VersionInfo struct {
	Major          int
	Minor          int
	Patch          int
	Build          int
	Branch         string
	Commit         string
	ReleaseChannel ReleaseChannel
}

func (v *VersionInfo) IsGreaterThan(b *VersionInfo) bool {
	if v.Major > b.Major {
		return true
	}

	if v.Major == b.Major &&
		v.Minor > b.Minor {
		return true
	}

	if v.Major == b.Major &&
		v.Minor == b.Minor &&
		v.Patch > b.Patch {
		return true
	}

	if v.Major == b.Major &&
		v.Minor == b.Minor &&
		v.Patch == b.Patch &&
		v.ReleaseChannel.GetPrio() > b.ReleaseChannel.GetPrio() {
		return true
	}

	if v.Major == b.Major &&
		v.Minor == b.Minor &&
		v.Patch == b.Patch &&
		v.ReleaseChannel.GetPrio() == b.ReleaseChannel.GetPrio() &&
		v.Build > b.Build {
		return true
	}

	return false
}

func (v *VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d.%d-%s [%s]", v.Major, v.Minor, v.Patch, v.Build, v.Branch, v.ReleaseChannel)
}
