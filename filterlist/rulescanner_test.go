package filterlist_test

import (
	"os"
	"strings"
	"testing"

	"github.com/AdguardTeam/urlfilter/filterlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleScanner_stringReader(t *testing.T) {
	t.Parallel()

	r := strings.NewReader(testRuleTextAll)
	scanner := filterlist.NewRuleScanner(r, filterListID, false)

	assert.True(t, scanner.Scan())
	f, idx := scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, testRule, f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
	assert.Equal(t, 0, idx)

	assert.True(t, scanner.Scan())
	f, idx = scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, testRuleCosmetic, f.Text())
	assert.Equal(t, filterListID, f.GetFilterListID())
	assert.Equal(t, cosmeticRuleIndex, idx)

	assert.False(t, scanner.Scan())
	assert.False(t, scanner.Scan())
}

func TestRuleScanner_fileReader(t *testing.T) {
	t.Parallel()

	file, err := os.Open(hostsPath)
	require.NoError(t, err)

	scanner := filterlist.NewRuleScanner(file, filterListID, true)
	rulesCount := 0
	for scanner.Scan() {
		f, idx := scanner.Rule()
		assert.NotNil(t, f)
		assert.Positive(t, idx)

		rulesCount++
	}

	assert.Equal(t, hostsRulesCount, rulesCount)
	assert.False(t, scanner.Scan())
}
