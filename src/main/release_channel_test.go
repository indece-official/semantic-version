package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseChannelEquals(t *testing.T) {
	a := ReleaseChannelAlpha
	b := ReleaseChannelAlpha
	c := ReleaseChannelFinal

	assert.True(t, a == b)
	assert.True(t, b == a)
	assert.False(t, a == c)
	assert.False(t, b == c)
}
