package main

type ReleaseChannel string

const (
	ReleaseChannelNone  ReleaseChannel = ""
	ReleaseChannelAlpha ReleaseChannel = "ALPHA"
	ReleaseChannelBeta  ReleaseChannel = "BETA"
	ReleaseChannelGamma ReleaseChannel = "GAMMA"
	ReleaseChannelFinal ReleaseChannel = "FINAL"
)

func (c ReleaseChannel) IsRelease() bool {
	return c != ReleaseChannelNone
}

func (c ReleaseChannel) GetPrio() int {
	switch c {
	case ReleaseChannelNone:
		return 0
	case ReleaseChannelAlpha:
		return 1
	case ReleaseChannelBeta:
		return 2
	case ReleaseChannelGamma:
		return 3
	case ReleaseChannelFinal:
		return 4
	default:
		return 0
	}
}
