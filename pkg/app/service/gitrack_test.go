package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	issueID1          = "ISSUE-1"
	issueID2          = "ISSUE-2"
	unexpectedIssueID = "UISSUE-2"
	issueTitle1       = "Нужно доработать класс экспорта"
	issueTitle2       = "Поправить e2e тесты"
	issueDescription1 = "В этой задаче нужно поправить Export Service чтобы он умел выводить даты. \n Сейчас он умеет только выводить строки"
	issueDescription2 = "Сейчас выпадает ошибка в e2e тестах err school not found\nНужно чтобы этот кейс обрабатывался согласно спеке"

	repository = "https://github.com/seelpy/repository.git"

	defaultBranch = "main"

	ytFeatureTag1        = "YTFeatureTag1"
	ytFeatureTagPre1     = "YTFeatureTagPre1"
	ytFeatureTagPublic1  = "YTFeatureTagPublic1"
	ytFeatureTagPost1    = "YTFeatureTagPublic1"
	featureBranchPre1    = "pre1"
	featureBranchPublic1 = "public1"
	featureBranchPost1   = "post1"
	nonExistentTag       = "non-existent"
)

var (
	issue1 = Issue{
		ID:          issueID1,
		Title:       issueTitle1,
		Description: issueDescription1,
		Tags: []string{
			ytFeatureTag1,
			ytFeatureTagPre1,
		},
	}
	issue2 = Issue{
		ID:          issueID2,
		Title:       issueTitle2,
		Description: issueDescription2,
		Tags: []string{
			ytFeatureTag1,
			ytFeatureTagPublic1,
		},
	}

	featureConfig1 = FeatureConfig{
		FeatureTag: ytFeatureTag1,
		Releases: []ReleaseConfig{
			{
				ReleaseBranch: featureBranchPre1,
				YoutrackTag:   ytFeatureTagPre1,
			},
			{
				ReleaseBranch: featureBranchPublic1,
				YoutrackTag:   ytFeatureTagPublic1,
			},
			{
				ReleaseBranch: featureBranchPost1,
				YoutrackTag:   ytFeatureTagPost1,
			},
		},
	}
)

func TestGetBranch(t *testing.T) {
	ctx := context.Background()

	git := newMockGit()
	youtrack := newMockYoutrack()
	configProvider := newMockGitrackConfigProvider()

	generateDefaultYoutrack(youtrack, git)
	generateDefaultConfigs(configProvider)

	gitrack := NewGitrack(git, youtrack, configProvider)

	git.checkout(issueID1)
	info, err := gitrack.GetBranchInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, issue1, info)

	git.checkout(issueID2)
	info, err = gitrack.GetBranchInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, issue2, info)
}

func TestGetBranchWithoutYoutrackIssue(t *testing.T) {
	ctx := context.Background()

	git := newMockGit()
	youtrack := newMockYoutrack()
	configProvider := newMockGitrackConfigProvider()

	generateDefaultYoutrack(youtrack, git)
	generateDefaultConfigs(configProvider)

	gitrack := NewGitrack(git, youtrack, configProvider)

	git.checkout(unexpectedIssueID)
	_, err := gitrack.GetBranchInfo(ctx)
	assert.ErrorIs(t, ErrIssueNotFound, err)
}

func TestGetBranchInDirectoryWithoutRepository(t *testing.T) {
	ctx := context.Background()

	git := newMockGit()
	youtrack := newMockYoutrack()
	configProvider := newMockGitrackConfigProvider()

	generateDefaultYoutrack(youtrack, git)
	generateDefaultConfigs(configProvider)
	git.setRepository("")

	gitrack := NewGitrack(git, youtrack, configProvider)

	git.checkout(issueID1)
	_, err := gitrack.GetBranchInfo(ctx)
	assert.ErrorIs(t, ErrGitRepositoryNotFound, err)
}

func TestCommit(t *testing.T) {
	git := newMockGit()
	youtrack := newMockYoutrack()
	configProvider := newMockGitrackConfigProvider()

	generateDefaultYoutrack(youtrack, git)
	generateDefaultConfigs(configProvider)

	git.checkout(issueID1)

	gitrack := NewGitrack(git, youtrack, configProvider)

	err := gitrack.Commit("init commit")
	assert.Nil(t, err)
	assert.Equal(t, "ISSUE-1 init commit", git.lastCommit)
}

func TestMerge(t *testing.T) {
	ctx := context.Background()

	git := newMockGit()
	youtrack := newMockYoutrack()
	configProvider := newMockGitrackConfigProvider()

	generateDefaultYoutrack(youtrack, git)
	generateDefaultConfigs(configProvider)

	git.checkout(issueID1)

	gitrack := NewGitrack(git, youtrack, configProvider)

	err := gitrack.Merge(ctx)
	assert.Nil(t, err)
	assert.Equal(t, []mergeData{
		{
			from: issueID1,
			to:   featureBranchPre1,
		},
		{
			from: featureBranchPre1,
			to:   featureBranchPublic1,
		},
		{
			from: featureBranchPublic1,
			to:   featureBranchPost1,
		},
	}, git.merges)

	git.clear()

	git.checkout(issueID2)
	err = gitrack.Merge(ctx)
	assert.Nil(t, err)
	assert.Equal(t, []mergeData{
		{
			from: issueID2,
			to:   featureBranchPublic1,
		},
		{
			from: featureBranchPublic1,
			to:   featureBranchPost1,
		},
	}, git.merges)
}

func generateDefaultYoutrack(youtrack *mockYoutrack, git *mockGit) {
	youtrack.addIssue(issue1)
	youtrack.addIssue(issue2)

	git.setRepository(repository)
}

func generateDefaultConfigs(configProvider *mockGitrackConfigProvider) {
	configProvider.addFeatureConfig(featureConfig1)
}

func newMockGit() *mockGit {
	return &mockGit{
		branch:     defaultBranch,
		repository: repository,
	}
}

func newMockYoutrack() *mockYoutrack {
	return &mockYoutrack{
		issues: make(map[string]Issue),
	}
}

func newMockGitrackConfigProvider() *mockGitrackConfigProvider {
	return &mockGitrackConfigProvider{
		configs: make(map[string]FeatureConfig),
	}
}

type mockGit struct {
	branch     string
	repository string
	lastCommit string
	merges     []mergeData
}

type mergeData struct {
	from string
	to   string
}

func (g *mockGit) GetBranch() (string, error) {
	if g.repository == "" {
		return "", ErrGitRepositoryNotFound
	}
	return g.branch, nil
}

func (g *mockGit) GetRepository() (string, error) {
	if g.repository == "" {
		return "", ErrGitRepositoryNotFound
	}
	return g.repository, nil
}

func (g *mockGit) Commit(message string) error {
	g.lastCommit = message
	return nil
}

func (g *mockGit) Merge(from string, to string) error {
	g.merges = append(g.merges, mergeData{
		from: from,
		to:   to,
	})
	return nil
}

func (g *mockGit) clear() {
	g.merges = make([]mergeData, 0)
}

func (g *mockGit) checkout(branch string) {
	g.branch = branch
}

func (g *mockGit) setRepository(repository string) {
	g.repository = repository
}

type mockYoutrack struct {
	issues map[string]Issue
}

func (y *mockYoutrack) GetIssue(_ context.Context, issueID string) (Issue, error) {
	if issue, ok := y.issues[issueID]; ok {
		return issue, nil
	}

	return Issue{}, ErrIssueNotFound
}

func (y *mockYoutrack) addIssue(issue Issue) {
	y.issues[issue.ID] = issue
}

type mockGitrackConfigProvider struct {
	configs map[string]FeatureConfig
}

func (m *mockGitrackConfigProvider) GetFeatureConfig(_ string, tag string) (FeatureConfig, error) {
	config, exists := m.configs[tag]
	if !exists {
		return FeatureConfig{}, ErrFeatureConfigNotFound
	}
	return config, nil
}

func (m *mockGitrackConfigProvider) addFeatureConfig(config FeatureConfig) {
	m.configs[config.FeatureTag] = config
}
