package filterlist_test

import (
	"strings"
	"testing"

	"github.com/AdguardTeam/urlfilter/filterlist"
	"github.com/stretchr/testify/assert"
)

func TestRuleStorageScanner(t *testing.T) {
	t.Parallel()

	const (
		filterList1 = testRuleTextAll
		filterList2 = "||example.com\n! test\n##advert"
	)

	r1 := strings.NewReader(filterList1)
	scanner1 := filterlist.NewRuleScanner(r1, 1, false)

	r2 := strings.NewReader(filterList2)
	scanner2 := filterlist.NewRuleScanner(r2, 2, false)

	// Now create the storage scanner.
	storageScanner := &filterlist.RuleStorageScanner{
		Scanners: []*filterlist.RuleScanner{scanner1, scanner2},
	}

	// Rule 1 from the list 1.
	assert.True(t, storageScanner.Scan())
	f, idx := storageScanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, testRule, f.Text())
	assert.Equal(t, 1, f.GetFilterListID())
	assert.Equal(t, "0x0000000100000000", int642hex(idx))

	// Rule 2 from the list 1.
	assert.True(t, storageScanner.Scan())
	f, idx = storageScanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, testRuleCosmetic, f.Text())
	assert.Equal(t, 1, f.GetFilterListID())
	assert.Equal(t, "0x0000000100000019", int642hex(idx))

	// Rule 1 from the list 2.
	assert.True(t, storageScanner.Scan())
	f, idx = storageScanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "||example.com", f.Text())
	assert.Equal(t, 2, f.GetFilterListID())
	assert.Equal(t, "0x0000000200000000", int642hex(idx))

	// Rule 2 from the list 2.
	assert.True(t, storageScanner.Scan())
	f, idx = storageScanner.Rule()

	assert.NotNil(t, f)
	assert.Equal(t, "##advert", f.Text())
	assert.Equal(t, 2, f.GetFilterListID())
	assert.Equal(t, "0x0000000200000015", int642hex(idx))

	// Now check that there's nothing to read.
	assert.False(t, storageScanner.Scan())

	// Check that nothing breaks if we read a finished scanner.
	assert.False(t, storageScanner.Scan())
}
