package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
)

type BranchPattern struct {
	exp *regexp.Regexp
}

func (p *BranchPattern) Match(str string) bool {
	return p.exp.MatchString(str)
}

func NewBranchPattern(pattern string) (*BranchPattern, error) {
	exp, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &BranchPattern{
		exp: exp,
	}, nil
}

func generateConfig() error {
	return GenerateConfig()
}

func getVersion() error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %s", err)
	}

	analyzer := NewAnalyzer(config)

	repo, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("error opening repository: %s", err)
	}

	err = analyzer.Load(repo)
	if err != nil {
		return fmt.Errorf("error loading analyzer: %s", err)
	}

	branchName, branchConfig, err := analyzer.GetCurrentBranchConfig(repo)
	if err != nil {
		return fmt.Errorf("error getting branch config: %s", err)
	}

	if branchConfig == nil {
		fmt.Printf("UNKNOWN\n")

		return nil
	}

	highestVersion, err := analyzer.GetHighestFinalReleaseVersion(repo)
	if err != nil {
		return fmt.Errorf("error getting highest final release: %s", err)
	}

	commits, err := analyzer.GetCommitsSinceLastRelease(repo, branchConfig, true)
	if err != nil {
		return fmt.Errorf("error loading commits since last release: %s", err)
	}

	commitParser := NewCommitParser()
	commitParser.Parse(commits)

	if highestVersion != nil {
		versionIncrement := commitParser.GetVersionIncrement()
		versionIncrement.Apply(highestVersion)
	} else {
		highestVersion = &VersionInfo{
			Major: 1,
			Minor: 0,
			Patch: 0,
			Build: 0,
		}
	}

	newTag, err := analyzer.GeneraterVersionTag(branchName, branchConfig, highestVersion)
	if err != nil {
		return fmt.Errorf("error generating version: %s", err)
	}

	// Output version
	fmt.Printf("%s\n", newTag)

	return nil
}

func getChangelog() error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %s", err)
	}

	analyzer := NewAnalyzer(config)

	repo, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("error opening repository: %s", err)
	}

	err = analyzer.Load(repo)
	if err != nil {
		return fmt.Errorf("error loading analyzer: %s", err)
	}

	_, branchConfig, err := analyzer.GetCurrentBranchConfig(repo)
	if err != nil {
		return fmt.Errorf("error getting branch config: %s", err)
	}

	if branchConfig == nil {
		fmt.Printf("\n")

		return nil
	}

	commits, err := analyzer.GetCommitsSinceLastRelease(repo, branchConfig, true)
	if err != nil {
		return fmt.Errorf("error loading commits since last release: %s", err)
	}

	commitParser := NewCommitParser()
	commitParser.Parse(commits)

	changelog := commitParser.GenerateChangelog()

	fmt.Printf("%s\n", changelog)

	return nil
}

func printHelp() {
	fmt.Printf("Usage: is2-semantic-version [args] <command>\n")
	fmt.Printf("\n")
	fmt.Printf("Args:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
	fmt.Printf("Commands:\n")
	fmt.Printf("  generate-config  Generate config file 'semanticversion.yaml'\n")
	fmt.Printf("  get-version      Get the new release version\n")
	fmt.Printf("  get-changelog    Get a changelog with all changes since the last release\n")
}

func main() {
	var err error

	flag.Parse()

	if len(flag.Args()) != 1 {
		printHelp()

		os.Exit(1)

		return
	}

	command := flag.Arg(0)

	switch command {
	case "generate-config":
		err = generateConfig()
	case "get-version":
		err = getVersion()
	case "get-changelog":
		err = getChangelog()
	default:
		printHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)

		os.Exit(1)

		return
	}
}
