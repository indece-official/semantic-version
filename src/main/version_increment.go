package main

type VersionIncrementLevel int

const (
	VersionIncrementLevelBuild VersionIncrementLevel = 0
	VersionIncrementLevelPatch VersionIncrementLevel = 1
	VersionIncrementLevelMinor VersionIncrementLevel = 2
	VersionIncrementLevelMajor VersionIncrementLevel = 3
)

type VersionIncrement struct {
	level VersionIncrementLevel
}

func (v *VersionIncrement) incrementTo(level VersionIncrementLevel) {
	if v.level >= level {
		return
	}

	v.level = level
}

func (v *VersionIncrement) IncrementMajor() {
	v.incrementTo(VersionIncrementLevelMajor)
}

func (v *VersionIncrement) IncrementMinor() {
	v.incrementTo(VersionIncrementLevelMinor)
}

func (v *VersionIncrement) IncrementPatch() {
	v.incrementTo(VersionIncrementLevelPatch)
}

func (v *VersionIncrement) IncrementBuild() {
	v.incrementTo(VersionIncrementLevelBuild)
}

func (v *VersionIncrement) Apply(versionInfo *VersionInfo) {
	switch v.level {
	case VersionIncrementLevelMajor:
		versionInfo.Major++
		versionInfo.Minor = 0
		versionInfo.Patch = 0
		versionInfo.Build = 0
	case VersionIncrementLevelMinor:
		versionInfo.Minor++
		versionInfo.Patch = 0
		versionInfo.Build = 0
	case VersionIncrementLevelPatch:
		versionInfo.Patch++
		versionInfo.Build = 0
	case VersionIncrementLevelBuild:
		versionInfo.Build++
	}
}

func NewVersionIncrement() *VersionIncrement {
	return &VersionIncrement{
		level: VersionIncrementLevelBuild,
	}
}
