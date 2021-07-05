package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ExpCleanBranchName = regexp.MustCompile(`[^a-zA-Z0-9\-\_]`)

type VersionPattern struct {
	releaseChannel ReleaseChannel
	pattern        string
	exp            *regexp.Regexp
}

func (p *VersionPattern) Parse(str string) *VersionInfo {
	versionInfo := &VersionInfo{}

	match := p.exp.FindStringSubmatch(str)
	if match == nil {
		return nil
	}

	versionInfo.ReleaseChannel = p.releaseChannel

	for i, name := range p.exp.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		switch name {
		case "major":
			major, err := strconv.Atoi(match[i])
			if err != nil {
				return nil
			}

			versionInfo.Major = major
		case "minor":
			minor, err := strconv.Atoi(match[i])
			if err != nil {
				return nil
			}

			versionInfo.Minor = minor
		case "patch":
			patch, err := strconv.Atoi(match[i])
			if err != nil {
				return nil
			}

			versionInfo.Patch = patch
		case "build":
			build, err := strconv.Atoi(match[i])
			if err != nil {
				return nil
			}

			versionInfo.Build = build
		case "branch":
			versionInfo.Branch = match[i]
		case "commit":
			versionInfo.Commit = match[i]
		case "shortcommit":
			versionInfo.ShortCommit = match[i]
		}
	}

	return versionInfo
}

func (v *VersionPattern) Generate(info *VersionInfo) string {
	branch := ExpCleanBranchName.ReplaceAllString(info.Branch, "_")

	str := v.pattern
	str = strings.ReplaceAll(str, "{major}", fmt.Sprintf("%d", info.Major))
	str = strings.ReplaceAll(str, "{minor}", fmt.Sprintf("%d", info.Minor))
	str = strings.ReplaceAll(str, "{patch}", fmt.Sprintf("%d", info.Patch))
	str = strings.ReplaceAll(str, "{build}", fmt.Sprintf("%d", info.Build))
	str = strings.ReplaceAll(str, "{branch}", branch)
	str = strings.ReplaceAll(str, "{commit}", info.Commit)
	str = strings.ReplaceAll(str, "{shortcommit}", info.ShortCommit)

	return str
}

func (v *VersionPattern) GenerateUnique(info *VersionInfo, usedTags map[string]bool, force bool) (string, error) {
	newTag := v.Generate(info)
	used, exists := usedTags[newTag]
	if !exists || !used {
		return newTag, nil
	}

	if !v.UsesBuild() {
		if force {
			return newTag, nil
		}

		// Iterating the build number doesn't make sens
		return "", fmt.Errorf("can't generate unique version without violating rules")
	}

	for {
		info.Build++
		newTag = v.Generate(info)
		used, exists = usedTags[newTag]
		if !exists || !used {
			return newTag, nil
		}
	}
}

func (v *VersionPattern) UsesBuild() bool {
	return strings.Contains(v.pattern, "{build}")
}

func NewVersionPattern(pattern string, releaseChannel ReleaseChannel) (*VersionPattern, error) {
	expPattern := pattern
	expPattern = strings.ReplaceAll(expPattern, "\\", "\\\\")
	expPattern = strings.ReplaceAll(expPattern, "-", "\\-")
	expPattern = strings.ReplaceAll(expPattern, ".", "\\.")
	expPattern = strings.ReplaceAll(expPattern, "{major}", "(?P<major>\\d+)")
	expPattern = strings.ReplaceAll(expPattern, "{minor}", "(?P<minor>\\d+)")
	expPattern = strings.ReplaceAll(expPattern, "{patch}", "(?P<patch>\\d+)")
	expPattern = strings.ReplaceAll(expPattern, "{build}", "(?P<build>\\d+)")
	expPattern = strings.ReplaceAll(expPattern, "{branch}", "(?P<branch>[a-zA-Z0-9\\_\\-\\\\/\\(\\)\\[\\]]+)")
	expPattern = strings.ReplaceAll(expPattern, "{commit}", "(?P<commit>[a-zA-Z0-9]+)")
	expPattern = strings.ReplaceAll(expPattern, "{shortcommit}", "(?P<shortcommit>[a-zA-Z0-9]+)")
	expPattern = fmt.Sprintf("^%s$", expPattern)

	exp, err := regexp.Compile(expPattern)
	if err != nil {
		return nil, err
	}

	return &VersionPattern{
		releaseChannel: releaseChannel,
		pattern:        pattern,
		exp:            exp,
	}, nil
}
