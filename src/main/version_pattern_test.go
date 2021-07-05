package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionPatternParse(t *testing.T) {
	ptr, err := NewVersionPattern("v{major}.{minor}.{patch}-{branch}.{build}", ReleaseChannelNone)
	assert.NoError(t, err)
	assert.NotNil(t, ptr)

	version := ptr.Parse("v1.23.12-feat/stephan.13")
	assert.NotNil(t, version)
	assert.Equal(t, 1, version.Major)
	assert.Equal(t, 23, version.Minor)
	assert.Equal(t, 12, version.Patch)
	assert.Equal(t, "feat/stephan", version.Branch)
	assert.Equal(t, "", version.Commit)
	assert.Equal(t, 13, version.Build)
}

func TestVersionGenerate(t *testing.T) {
	ptr, err := NewVersionPattern("v{major}.{minor}.{patch}-{branch}.{build}", ReleaseChannelNone)
	assert.NoError(t, err)
	assert.NotNil(t, ptr)

	versionInfo := &VersionInfo{
		Major:  2,
		Minor:  12,
		Patch:  56,
		Build:  0,
		Branch: "feat/stephan",
		Commit: "abcdef12345",
	}

	version := ptr.Generate(versionInfo)
	assert.Equal(t, "v2.12.56-feat_stephan.0", version)
}

func TestVersionGenerateShortCommit(t *testing.T) {
	ptr, err := NewVersionPattern("v{major}.{minor}.{patch}-{shortcommit}.{build}", ReleaseChannelNone)
	assert.NoError(t, err)
	assert.NotNil(t, ptr)

	versionInfo := &VersionInfo{
		Major:       2,
		Minor:       12,
		Patch:       56,
		Build:       0,
		Branch:      "feat/stephan",
		ShortCommit: "abcdef12345",
	}

	version := ptr.Generate(versionInfo)
	assert.Equal(t, "v2.12.56-abcdef12345.0", version)
}

func TestVersionGenerateUnique(t *testing.T) {
	ptr, err := NewVersionPattern("v{major}.{minor}.{patch}-{branch}.{build}", ReleaseChannelNone)
	assert.NoError(t, err)
	assert.NotNil(t, ptr)

	versionInfo := &VersionInfo{
		Major:  2,
		Minor:  12,
		Patch:  56,
		Build:  0,
		Branch: "feat/stephan",
		Commit: "abcdef12345",
	}

	usedTags := map[string]bool{
		"v2.12.56":                true,
		"v2.12.56-feat_stephan.0": true,
		"v2.12.56-feat_stephan.1": true,
		"v2.12.56-feat_stephan.2": true,
	}

	version, err := ptr.GenerateUnique(versionInfo, usedTags, false)
	assert.NoError(t, err)
	assert.Equal(t, "v2.12.56-feat_stephan.3", version)
}

func TestVersionGenerateUniqueNotPossible(t *testing.T) {
	ptr, err := NewVersionPattern("v{major}.{minor}.{patch}", ReleaseChannelFinal)
	assert.NoError(t, err)
	assert.NotNil(t, ptr)

	versionInfo := &VersionInfo{
		Major:  2,
		Minor:  12,
		Patch:  56,
		Build:  0,
		Branch: "master",
		Commit: "abcdef12345",
	}

	usedTags := map[string]bool{
		"v2.12.56":                true,
		"v2.12.56-feat_stephan.0": true,
		"v2.12.56-feat_stephan.1": true,
		"v2.12.56-feat_stephan.2": true,
	}

	_, err = ptr.GenerateUnique(versionInfo, usedTags, false)
	assert.Error(t, err)
}

func TestVersionGenerateUniqueNotPossibleForce(t *testing.T) {
	ptr, err := NewVersionPattern("v{major}.{minor}.{patch}", ReleaseChannelFinal)
	assert.NoError(t, err)
	assert.NotNil(t, ptr)

	versionInfo := &VersionInfo{
		Major:  2,
		Minor:  12,
		Patch:  56,
		Build:  0,
		Branch: "master",
		Commit: "abcdef12345",
	}

	usedTags := map[string]bool{
		"v2.12.56":                true,
		"v2.12.56-feat_stephan.0": true,
		"v2.12.56-feat_stephan.1": true,
		"v2.12.56-feat_stephan.2": true,
	}

	version, err := ptr.GenerateUnique(versionInfo, usedTags, true)
	assert.NoError(t, err)
	assert.Equal(t, "v2.12.56", version)
}
