package filterlist_test

import (
	"fmt"
	"strings"
)

// filterListID is the testing ID for a filter list.
const filterListID = 1

// Common domains for tests.
const testDomain = "test.example"

// Common rules for tests.
const (
	testRule         = "||" + testDomain
	testRuleCosmetic = "##banner"
	testComment      = "! comment"
)

// Common text rules for tests.
const (
	testRuleText         = testRule + "\n"
	testRuleTextCosmetic = testRuleCosmetic + "\n"
	testCommentText      = testComment + "\n"

	testRuleTextAll = testRuleText + testCommentText + testRuleTextCosmetic
)

// cosmeticRuleIndex is the index of the cosmetic rule in testRuleTextAll.
var cosmeticRuleIndex = strings.Index(testRuleTextAll, testRuleCosmetic)

const (
	// testResourcesDir is the path to test resources.
	testResourcesDir = "../testdata"

	// hostsPath is the path to hosts file for testing.
	hostsPath = testResourcesDir + "/hosts"

	// hostsRulesCount is the number of rules in the hosts file available by
	// hostsPath.
	//
	// NOTE:  Keep in sync with hostsPath file contents.
	hostsRulesCount = 55997
)

// int642hex returns a hex string representation of an int64 value.
//
// TODO(d.kolyshev):  Remove.
func int642hex(v int64) (s string) {
	return fmt.Sprintf("0x%016x", v)
}
