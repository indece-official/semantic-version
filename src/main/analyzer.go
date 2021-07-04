package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Analyzer struct {
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

	a.headCommit, err = repo.CommitObject(a.head.Hash())
	if err != nil {
		return fmt.Errorf("can't load head commit: %s", err)
	}

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
	branchIter, err := repo.Branches()
	if err != nil {
		return "", nil, fmt.Errorf("can't load branches: %s", err)
	}

	branchName := ""

	for {
		branch, err := branchIter.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return "", nil, fmt.Errorf("can't iterate branches: %s", err)
		}

		if branch.Hash() == a.head.Hash() {
			branchName = branch.Name().Short()
			break
		}
	}

	if branchName == "" {
		return "", nil, nil
	}

	for _, branchConfig := range a.config.Branches {
		if branchConfig.GetBranchPattern().Match(branchName) {
			return branchName, branchConfig, nil
		}
	}

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

	return highestTag, nil
}

func (a *Analyzer) GetCommitsSinceLastRelease(repo *git.Repository, branchConfig *BranchConfig, channelSensitive bool) ([]*object.Commit, error) {
	commits := []*object.Commit{}

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
				if !versionInfo.ReleaseChannel.IsRelease() {
					continue
				}

				if !channelSensitive || versionInfo.ReleaseChannel == branchConfig.ReleaseChannel {
					// Found matching release commit, return
					return commits, nil
				}
			}
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func (a *Analyzer) GeneraterVersionTag(branchName string, branchConfig *BranchConfig, versionInfo *VersionInfo) (string, error) {
	versionInfo.Branch = branchName
	versionInfo.Commit = a.headCommit.Hash.String()

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
