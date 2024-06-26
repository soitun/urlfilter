package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMatchingResult(t *testing.T) {
	rules := testNewNetworkRules(t, []string{
		"||example.org^",
	}, 0)
	sourceRules := []*NetworkRule{}
	result := NewMatchingResult(rules, sourceRules)

	assert.NotNil(t, result.BasicRule)
	assert.NotNil(t, result.GetBasicResult())
	assert.Equal(t, "||example.org^", result.GetBasicResult().String())
}

func TestNewMatchingResultWhitelist(t *testing.T) {
	rules := testNewNetworkRules(t, []string{
		"||example.org^",
	}, 0)
	sourceRules := testNewNetworkRules(t, []string{
		"@@||example.com^$document",
	}, 0)
	result := NewMatchingResult(rules, sourceRules)

	assert.Nil(t, result.BasicRule)
	assert.NotNil(t, result.DocumentRule)
	assert.NotNil(t, result.GetBasicResult())
	assert.Equal(t, "@@||example.com^$document", result.GetBasicResult().String())
}

func TestGetCosmeticOption(t *testing.T) {
	// Simple case - no limitations
	rules := testNewNetworkRules(t, []string{
		"||example.org^",
	}, 0)
	sourceRules := []*NetworkRule{}
	result := NewMatchingResult(rules, sourceRules)
	assert.Equal(t, CosmeticOptionAll, result.GetCosmeticOption())

	// $generichide
	rules = testNewNetworkRules(t, []string{
		"@@||example.org^$generichide",
	}, 0)
	sourceRules = []*NetworkRule{}
	result = NewMatchingResult(rules, sourceRules)
	assert.Equal(t, CosmeticOptionCSS|CosmeticOptionJS, result.GetCosmeticOption())

	// $jsinject
	rules = testNewNetworkRules(t, []string{
		"@@||example.org^$jsinject",
	}, 0)
	sourceRules = []*NetworkRule{}
	result = NewMatchingResult(rules, sourceRules)
	assert.Equal(t, CosmeticOptionCSS|CosmeticOptionGenericCSS, result.GetCosmeticOption())

	// $elemhide
	rules = testNewNetworkRules(t, []string{
		"@@||example.org^$elemhide",
	}, 0)
	sourceRules = []*NetworkRule{}
	result = NewMatchingResult(rules, sourceRules)
	assert.Equal(t, CosmeticOptionJS, result.GetCosmeticOption())

	// $document
	rules = testNewNetworkRules(t, []string{
		"@@||example.org^$document",
	}, 0)
	sourceRules = []*NetworkRule{}
	result = NewMatchingResult(rules, sourceRules)
	assert.Equal(t, CosmeticOptionNone, result.GetCosmeticOption())
}

func TestNewMatchingResultBadfilter(t *testing.T) {
	rules := testNewNetworkRules(t, []string{
		"||example.org^",
		"||example.org^$badfilter",
	}, 0)
	sourceRules := []*NetworkRule{}
	result := NewMatchingResult(rules, sourceRules)

	assert.Nil(t, result.BasicRule)
	assert.Nil(t, result.DocumentRule)
}

func TestNewMatchingResultBadfilterWhitelist(t *testing.T) {
	rules := testNewNetworkRules(t, []string{
		"||example.org^",
		"@@||example.org^",
		"@@||example.org^$badfilter",
	}, 0)
	sourceRules := []*NetworkRule{}
	result := NewMatchingResult(rules, sourceRules)

	assert.NotNil(t, result.BasicRule)
	assert.Nil(t, result.DocumentRule)
	assert.Equal(t, "||example.org^", result.GetBasicResult().String())
}

func TestNewMatchingResultBadfilterSourceRules(t *testing.T) {
	rules := testNewNetworkRules(t, []string{
		"||example.org^",
	}, 0)
	sourceRules := testNewNetworkRules(t, []string{
		"@@||example.org^$document",
		"@@||example.org^$document,badfilter",
	}, 0)
	result := NewMatchingResult(rules, sourceRules)

	assert.NotNil(t, result.BasicRule)
	assert.Nil(t, result.DocumentRule)
	assert.Equal(t, "||example.org^", result.GetBasicResult().String())
}

// TODO: ADD MORE TESTS

// testNewNetworkRules creates a list of network rules from a string array
func testNewNetworkRules(t *testing.T, lines []string, filterListID int) []*NetworkRule {
	var rules []*NetworkRule

	for _, line := range lines {
		f, err := NewNetworkRule(line, filterListID)
		if err != nil {
			t.Fatalf("failed to create network rule from %s: %s", line, err)
		}
		rules = append(rules, f)
	}

	return rules
}

func TestRemoveDNSRewriteRules(t *testing.T) {
	rules := []*NetworkRule{{
		RuleText:   "host1",
		DNSRewrite: nil,
	}, {
		RuleText:   "host2",
		DNSRewrite: nil,
	}, {
		RuleText:   "host3",
		DNSRewrite: nil,
	}}

	got := removeDNSRewriteRules(rules)
	assert.Equal(t, rules, got)

	rules = []*NetworkRule{{
		RuleText:   "host1",
		DNSRewrite: nil,
	}, {
		RuleText:   "host2",
		DNSRewrite: &DNSRewrite{},
	}, {
		RuleText:   "host3",
		DNSRewrite: nil,
	}}

	got = removeDNSRewriteRules(rules)
	assert.NotEqual(t, rules, got)
	if assert.Equal(t, 2, len(got)) {
		assert.Equal(t, "host1", got[0].RuleText)
		assert.Equal(t, "host3", got[1].RuleText)
	}
}

func TestGetDNSBasicRule(t *testing.T) {
	blockRule := newTestRule(t, "example.block")
	allowlistRule := newTestRule(t, "@@||example.allow^")
	importantBlockRule := newTestRule(t, "example.block$important")

	testCases := []struct {
		want  *NetworkRule
		name  string
		rules []*NetworkRule
	}{{
		want:  nil,
		rules: []*NetworkRule{},
		name:  "empty",
	}, {
		want: blockRule,
		rules: []*NetworkRule{
			blockRule,
		},
		name: "basic",
	}, {
		want: allowlistRule,
		rules: []*NetworkRule{
			blockRule,
			allowlistRule,
		},
		name: "allowlist",
	}, {
		want: importantBlockRule,
		rules: []*NetworkRule{
			blockRule,
			allowlistRule,
			importantBlockRule,
		},
		name: "important",
	}, {
		want: blockRule,
		rules: []*NetworkRule{
			blockRule,
			newTestRule(t, "@@||example.org^$stealth"),
		},
		name: "stealth",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := GetDNSBasicRule(tc.rules)
			assert.Equal(t, tc.want, r)
		})
	}
}

// newTestRule returns a network rule created from given source text.
func newTestRule(t *testing.T, srcText string) (r *NetworkRule) {
	t.Helper()

	r, err := NewNetworkRule(srcText, 1)
	require.NoError(t, err)

	return r
}
