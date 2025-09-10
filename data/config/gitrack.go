package config

import (
	"github.com/pkg/errors"
	appservice "gitrack/pkg/app/service"
	"gopkg.in/yaml.v3"
	"os"
)

type GitrackConfigProvider interface {
	GetFeatureConfig(tag string, repository string) (FeatureConfig, error)
}

type Features []FeatureConfig

type FeatureConfig struct {
	FeatureTag string
	Releases   []ReleaseConfig
}

type ReleaseConfig struct {
	ReleaseBranch string
	YoutrackTag   string
}

type Config struct {
	Gitrack GitrackConfig `yaml:"gitrack"`
}

type GitrackConfig struct {
	Youtrack YoutrackConfig `yaml:"youtrack"`
	Features []Feature      `yaml:"features"`
}

type YoutrackConfig struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

type Feature struct {
	BaseTag      string       `yaml:"baseTag"`
	Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
	RepositoryName string    `yaml:"repository"`
	Releases       []Release `yaml:"releases"`
}

type Release struct {
	Tag    string `yaml:"tag"`
	Branch string `yaml:"branch"`
}

func ParseConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}

func (c *Config) GetFeatureConfig(tag string, repository string) (FeatureConfig, error) {
	for _, feature := range c.Gitrack.Features {
		for _, repo := range feature.Repositories {
			if repo.RepositoryName == repository {
				for _, release := range repo.Releases {
					if release.Tag == tag {
						return FeatureConfig{
							FeatureTag: feature.BaseTag,
							Releases:   convertReleases(repo.Releases),
						}, nil
					}
				}
			}
		}
	}

	return FeatureConfig{}, appservice.ErrFeatureConfigNotFound
}

func convertReleases(releases []Release) []ReleaseConfig {
	var result []ReleaseConfig
	for _, release := range releases {
		result = append(result, ReleaseConfig{
			ReleaseBranch: release.Branch,
			YoutrackTag:   release.Tag,
		})
	}
	return result
}
