package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionIncrementApply(t *testing.T) {
	info := &VersionInfo{
		Major: 1,
		Minor: 0,
		Patch: 2,
		Build: 0,
	}

	inc := NewVersionIncrement()
	inc.IncrementBuild()
	inc.IncrementMinor()
	inc.Apply(info)

	assert.Equal(t, 1, info.Major)
	assert.Equal(t, 1, info.Minor)
	assert.Equal(t, 0, info.Patch)
	assert.Equal(t, 0, info.Build)
}
