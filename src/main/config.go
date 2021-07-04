package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigFilename = "./semanticversion.yaml"

type BranchConfig struct {
	BranchPattern string `yaml:"branch_pattern"`

	// VersionPattern specifies the pattern to generate new version numbers
	// Placeholders:
	//   {major} Major version
	//   {minor} Minor version
	//   {path} Path version
	//   {branch} Branch name
	//   {commit} Commit hash (short)
	//   {build} Build number
	VersionPattern string         `yaml:"version_pattern"`
	ReleaseChannel ReleaseChannel `yaml:"release_channel"`

	branchPattern  *BranchPattern
	versionPattern *VersionPattern
}

func (c *BranchConfig) GetBranchPattern() *BranchPattern {
	return c.branchPattern
}

func (c *BranchConfig) GetVersionPattern() *VersionPattern {
	return c.versionPattern
}

func (c *BranchConfig) Parse() error {
	var err error

	c.branchPattern, err = NewBranchPattern(c.BranchPattern)
	if err != nil {
		return fmt.Errorf("can't parse branch pattern \"%s\": %s", c.BranchPattern, err)
	}

	c.versionPattern, err = NewVersionPattern(c.VersionPattern, c.ReleaseChannel)
	if err != nil {
		return fmt.Errorf("can't parse version pattern \"%s\": %s", c.VersionPattern, err)
	}

	return nil
}

type Config struct {
	Branches []*BranchConfig `yaml:"branches"`
}

func (c *Config) Parse() error {
	for _, branch := range c.Branches {
		err := branch.Parse()
		if err != nil {
			return err
		}
	}

	return nil
}

var DefaultConfig = &Config{
	Branches: []*BranchConfig{
		{
			BranchPattern:  "master",
			VersionPattern: "v{major}.{minor}.{patch}",
			ReleaseChannel: ReleaseChannelFinal,
		},
		{
			BranchPattern:  "release.*",
			VersionPattern: "v{major}.{minor}.{patch}",
			ReleaseChannel: ReleaseChannelFinal,
		},
		{
			BranchPattern:  "gamma.*",
			VersionPattern: "v{major}.{minor}.{patch}-gamma.{build}",
			ReleaseChannel: ReleaseChannelGamma,
		},
		{
			BranchPattern:  "beta.*",
			VersionPattern: "v{major}.{minor}.{patch}-beta.{build}",
			ReleaseChannel: ReleaseChannelBeta,
		},
		{
			BranchPattern:  "alpha.*",
			VersionPattern: "v{major}.{minor}.{patch}-alpha.{build}",
			ReleaseChannel: ReleaseChannelAlpha,
		},
		{
			BranchPattern:  "feat.*",
			VersionPattern: "v{major}.{minor}.{patch}-{branch}.{build}",
		},
		{
			BranchPattern:  "fix.*",
			VersionPattern: "v{major}.{minor}.{patch}-{branch}.{build}",
		},
	},
}

func GenerateConfig() error {
	_, err := os.Stat(ConfigFilename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("can't access file %s: %s", ConfigFilename, err)
	}

	if err == nil {
		return fmt.Errorf("file %s already exists: %s", ConfigFilename, err)
	}

	configData, err := yaml.Marshal(DefaultConfig)
	if err != nil {
		return fmt.Errorf("can't encode default config: %s", err)
	}

	err = ioutil.WriteFile(ConfigFilename, configData, 0)
	if err != nil {
		return fmt.Errorf("can't write default config to file %s: %s", ConfigFilename, err)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	_, err := os.Stat(ConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			err = DefaultConfig.Parse()
			if err != nil {
				return nil, err
			}

			return DefaultConfig, nil
		}

		return nil, fmt.Errorf("can't open config file %s: %s", ConfigFilename, err)
	}

	config := &Config{}

	configData, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		return nil, fmt.Errorf("can't read config file %s: %s", ConfigFilename, err)
	}

	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return nil, fmt.Errorf("can't parse config file %s: %s", ConfigFilename, err)
	}

	err = config.Parse()
	if err != nil {
		return nil, err
	}

	return config, nil
}
