package main

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var flagGitBranch = flag.String("git-branch", "", "")
var flagBuild = flag.Int("build", -1, "")

type Analyzer struct {
	// commit-hash => VersionInfo of tag
	mapCommitTags map[string][]*VersionInfo
	mapTags       map[string]bool
	head          *plumbing.Reference
	headCommit    *object.Commit
	config        *Config
}

func (a *Analyzer) Load(repo *git.Repository) error {
	var err error

	a.head, err = repo.Head()
	if err != nil {
		return fmt.Errorf("can't load head: %s", err)
	}

	Debugf("Head is %s", a.head.Hash().String())

	a.headCommit, err = repo.CommitObject(a.head.Hash())
	if err != nil {
		return fmt.Errorf("can't load head commit: %s", err)
	}

	Debugf("Head commit is %s", a.headCommit.Hash.String())

	tags, err := repo.Tags()
	if err != nil {
		return fmt.Errorf("can't load tags: %s", err)
	}

	for {
		tag, err := tags.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("can't iterate tags: %s", err)
		}

		tagName := tag.Name().Short()
		revision := plumbing.Revision(tagName)
		tagCommit, err := repo.ResolveRevision(revision)
		if err != nil {
			return fmt.Errorf("can't resolve tag %s: %s", tagName, err)
		}
		tagCommitStr := tagCommit.String()

		for _, branchConfig := range a.config.Branches {
			versionInfo := branchConfig.GetVersionPattern().Parse(tagName)
			if versionInfo == nil {
				continue
			}

			Debugf("Found tag %s (%s) => %v", tagName, tagCommitStr, versionInfo)

			if _, exists := a.mapCommitTags[tagCommitStr]; !exists {
				a.mapCommitTags[tagCommitStr] = []*VersionInfo{}
			}

			a.mapCommitTags[tagCommitStr] = append(a.mapCommitTags[tagCommitStr], versionInfo)
			a.mapTags[tagName] = true
		}
	}

	return nil
}

func (a *Analyzer) GetCurrentBranchConfig(repo *git.Repository) (string, *BranchConfig, error) {
	branchName := ""

	if *flagGitBranch != "" {
		branchName = *flagGitBranch
	} else {
		branchName = a.head.Name().Short()
	}

	if branchName == "" || branchName == "HEAD" {
		Debugf("Found no valid branch name: %s", branchName)

		return "", nil, nil
	}

	for _, branchConfig := range a.config.Branches {
		if branchConfig.GetBranchPattern().Match(branchName) {
			Debugf("Found config %s for branch name %s", branchConfig.BranchPattern, branchName)

			return branchName, branchConfig, nil
		}
	}

	Debugf("Found no config for branch name %s", branchName)

	return branchName, nil, nil
}

func (a *Analyzer) GetHighestFinalReleaseVersion(repo *git.Repository) (*VersionInfo, error) {
	var highestTag *VersionInfo

	commitIter := object.NewCommitPostorderIter(a.headCommit, []plumbing.Hash{})
	for {
		commit, err := commitIter.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("can't iterate commits: %s", err)
		}

		versionInfos, exists := a.mapCommitTags[commit.Hash.String()]
		if exists {
			for _, versionInfo := range versionInfos {
				if versionInfo.ReleaseChannel != ReleaseChannelFinal {
					continue
				}

				if highestTag == nil || versionInfo.IsGreaterThan(highestTag) {
					highestTag = versionInfo
				}
			}
		}
	}

	Debugf("Found highest release tag %v", highestTag)

	return highestTag, nil
}

func (a *Analyzer) GetCommitsSinceLastRelease(repo *git.Repository, branchConfig *BranchConfig, channelSensitive bool) ([]*object.Commit, error) {
	commits := []*object.Commit{}

	seenExternal := map[plumbing.Hash]bool{}
	for commitHash, versionInfos := range a.mapCommitTags {
		for _, versionInfo := range versionInfos {
			if !versionInfo.ReleaseChannel.IsRelease() {
				continue
			}

			if !channelSensitive || versionInfo.ReleaseChannel.GetPrio() >= branchConfig.ReleaseChannel.GetPrio() {
				// Found matching release commit
				commit, err := repo.CommitObject(plumbing.NewHash(commitHash))
				if err != nil {
					return nil, fmt.Errorf("can't load commit object: %s", err)
				}

				commitIter := object.NewCommitPreorderIter(commit, seenExternal, []plumbing.Hash{})
				commitIter.ForEach(func(c *object.Commit) error {
					seenExternal[c.Hash] = true

					return nil
				})

				seenExternal[plumbing.NewHash(commitHash)] = true

				break
			}
		}
	}

	commitIter := object.NewCommitIterBSF(a.headCommit, seenExternal, []plumbing.Hash{})
	for {
		commit, err := commitIter.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("can't iterate commits: %s", err)
		}

		Debugf("Analyze commit %s for changelog => %v", commit.Hash.String(), a.mapCommitTags[commit.Hash.String()])

		commits = append(commits, commit)
	}

	return commits, nil
}

func (a *Analyzer) GeneraterVersionTag(branchName string, branchConfig *BranchConfig, versionInfo *VersionInfo) (string, error) {
	versionInfo.Branch = branchName
	versionInfo.Commit = a.headCommit.Hash.String()
	versionInfo.ShortCommit = a.headCommit.Hash.String()[:10]

	if *flagBuild >= 0 {
		versionInfo.Build = *flagBuild
	}

	newTag, err := branchConfig.GetVersionPattern().GenerateUnique(versionInfo, a.mapTags, true)
	if err != nil {
		return "", err
	}

	return newTag, nil
}

func NewAnalyzer(config *Config) *Analyzer {
	return &Analyzer{
		mapCommitTags: map[string][]*VersionInfo{},
		mapTags:       map[string]bool{},
		config:        config,
	}
}
