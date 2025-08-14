package filterlist_test

import (
	"path/filepath"
	"testing"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/urlfilter/filterlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO(d.kolyshev):  Improve tests.

func TestStringRuleListScanner(t *testing.T) {
	t.Parallel()

	ruleList := &filterlist.StringRuleList{
		ID:             filterListID,
		IgnoreCosmetic: false,
		RulesText:      "||example.org\n! test\n##banner",
	}
	testutil.CleanupAndRequireSuccess(t, ruleList.Close)

	assert.Equal(t, 1, ruleList.GetID())

	scanner := ruleList.NewScanner()

	assert.True(t, scanner.Scan())
	f, idx := scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "||example.org", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
	assert.Equal(t, 0, idx)

	assert.True(t, scanner.Scan())
	f, idx = scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "##banner", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
	assert.Equal(t, 21, idx)

	// Finish scanning
	assert.False(t, scanner.Scan())

	f, err := ruleList.RetrieveRule(0)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "||example.org", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())

	f, err = ruleList.RetrieveRule(21)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "##banner", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
}

func TestFileRuleListScanner(t *testing.T) {
	t.Parallel()

	testFileRuleList := filepath.Join(testResourcesDir, "test_file_rule_list.txt")
	ruleList, err := filterlist.NewFileRuleList(filterListID, testFileRuleList, false)
	require.NoError(t, err)
	testutil.CleanupAndRequireSuccess(t, ruleList.Close)

	assert.Equal(t, filterListID, ruleList.GetID())

	scanner := ruleList.NewScanner()

	assert.True(t, scanner.Scan())
	f, idx := scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "||example.org", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
	assert.Equal(t, 0, idx)

	assert.True(t, scanner.Scan())
	f, idx = scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "##banner", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
	assert.Equal(t, 21, idx)

	// Finish scanning
	assert.False(t, scanner.Scan())

	f, err = ruleList.RetrieveRule(0)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "||example.org", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())

	f, err = ruleList.RetrieveRule(21)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "##banner", f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
}
