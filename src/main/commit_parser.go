package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

var expBreak = regexp.MustCompile("^[bB][rR][eE][aA][kK]:")
var expFeat = regexp.MustCompile("^[fF][eE][aA][tT]:")
var expFix = regexp.MustCompile("^[fF][iI][xX]:")

type ParsedCommit struct {
	Message string
	Hash    string
}

type CommitParser struct {
	Major   []*ParsedCommit
	Minor   []*ParsedCommit
	Patch   []*ParsedCommit
	Unknown []*ParsedCommit
}

func (c *CommitParser) Parse(commits []*object.Commit) {
	for _, commit := range commits {
		commitParts := strings.Split(commit.Message, ";")
		for _, commitPart := range commitParts {
			commitPart = strings.TrimSpace(commitPart)
			commitPartMessage := ""
			isMajor := false
			isMinor := false
			isPatch := false

			switch {
			case expBreak.MatchString(commitPart):
				isMajor = true
				commitPartMessage = expBreak.ReplaceAllString(commitPart, "")
			case expFeat.MatchString(commitPart):
				isMinor = true
				commitPartMessage = expFeat.ReplaceAllString(commitPart, "")
			case expFix.MatchString(commitPart):
				isPatch = true
				commitPartMessage = expFix.ReplaceAllString(commitPart, "")
			default:
				commitPartMessage = commitPart
			}

			commitPartMessage = strings.TrimSpace(commitPartMessage)
			parsedCommit := &ParsedCommit{
				Message: commitPartMessage,
				Hash:    commit.Hash.String(),
			}

			switch {
			case isMajor:
				c.Major = append(c.Major, parsedCommit)
			case isMinor:
				c.Minor = append(c.Minor, parsedCommit)
			case isPatch:
				c.Patch = append(c.Patch, parsedCommit)
			default:
				c.Unknown = append(c.Unknown, parsedCommit)
			}
		}
	}
}

func (c *CommitParser) GenerateChangelog() string {
	msg := ""

	if len(c.Major) > 0 {
		msg += "# BREAKING CHANGES\n"
		for _, parseCommit := range c.Major {
			msg += fmt.Sprintf("* %s (%s)\n", parseCommit.Message, parseCommit.Hash)
		}
		msg += "\n"
	}

	if len(c.Minor) > 0 {
		msg += "# Features\n"
		for _, parseCommit := range c.Minor {
			msg += fmt.Sprintf("* %s (%s)\n", parseCommit.Message, parseCommit.Hash)
		}
		msg += "\n"
	}

	if len(c.Patch) > 0 {
		msg += "# Fixes\n"
		for _, parseCommit := range c.Patch {
			msg += fmt.Sprintf("* %s (%s)\n", parseCommit.Message, parseCommit.Hash)
		}
		msg += "\n"
	}

	return msg
}

func (c *CommitParser) GetVersionIncrement() *VersionIncrement {
	versionIncrement := NewVersionIncrement()
	switch {
	case len(c.Major) > 0:
		versionIncrement.IncrementMajor()
	case len(c.Minor) > 0:
		versionIncrement.IncrementMinor()
	case len(c.Patch) > 0:
		versionIncrement.IncrementPatch()
	case len(c.Unknown) > 0:
		versionIncrement.IncrementBuild()
	}

	return versionIncrement
}

func NewCommitParser() *CommitParser {
	return &CommitParser{}
}
