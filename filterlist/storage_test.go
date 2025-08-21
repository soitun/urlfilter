package filterlist_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/urlfilter/filterlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleStorage(t *testing.T) {
	t.Parallel()

	list1 := filterlist.NewString(&filterlist.StringConfig{
		RulesText: testRuleTextAll,
		ID:        1,
	})

	list2 := filterlist.NewString(&filterlist.StringConfig{
		RulesText: "||example.com\n! test\n##advert",
		ID:        2,
	})

	// Create storage from two lists.
	storage, err := filterlist.NewRuleStorage([]filterlist.Interface{list1, list2})
	require.NoError(t, err)

	// Create a scanner instance.
	scanner := storage.NewRuleStorageScanner()
	assert.NotNil(t, scanner)

	// Time to scan!

	// Rule 1 from the list 1
	assert.True(t, scanner.Scan())
	f, idx := scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, testRule, f.Text())
	assert.Equal(t, 1, f.GetFilterListID())
	assert.Equal(t, "0x0000000100000000", int642hex(idx))

	// Rule 2 from the list 1.
	assert.True(t, scanner.Scan())
	f, idx = scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, testRuleCosmetic, f.Text())
	assert.Equal(t, 1, f.GetFilterListID())
	assert.Equal(t, "0x0000000100000019", int642hex(idx))

	// Rule 1 from the list 2.
	assert.True(t, scanner.Scan())
	f, idx = scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "||example.com", f.Text())
	assert.Equal(t, 2, f.GetFilterListID())
	assert.Equal(t, "0x0000000200000000", int642hex(idx))

	// Rule 2 from the list 2.
	assert.True(t, scanner.Scan())
	f, idx = scanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "##advert", f.Text())
	assert.Equal(t, 2, f.GetFilterListID())
	assert.Equal(t, "0x0000000200000015", int642hex(idx))

	// Now check that there's nothing to read.
	assert.False(t, scanner.Scan())

	// Check that nothing breaks if we read a finished scanner.
	assert.False(t, scanner.Scan())

	// Time to retrieve!

	// Rule 1 from the list 1.
	f, err = storage.RetrieveRule(0x0000000100000000)
	require.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, testRule, f.Text())

	// Rule 2 from the list 1.
	f, err = storage.RetrieveRule(0x0000000100000019)
	require.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, testRuleCosmetic, f.Text())

	// Rule 1 from the list 2.
	f, err = storage.RetrieveRule(0x0000000200000000)
	require.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "||example.com", f.Text())

	// Rule 2 from the list 2.
	f, err = storage.RetrieveRule(0x0000000200000015)
	require.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "##advert", f.Text())
}

func TestRuleStorage_invalid(t *testing.T) {
	t.Parallel()

	conf := &filterlist.StringConfig{
		ID:        filterListID,
		RulesText: "",
	}
	_, err := filterlist.NewRuleStorage([]filterlist.Interface{
		filterlist.NewString(conf),
		filterlist.NewString(conf),
	})
	testutil.AssertErrorMsg(t, "at index 1: id: duplicated value: '\\x01'", err)
}
