// Package lookup implements index structures that we use to improve matching
// speed in the engines.
package lookup

import "github.com/AdguardTeam/urlfilter/rules"

// Table is the interface for all lookup tables used to speed up matching.
type Table interface {
	// Add adds the rule to the lookup table.  If ok is false, the rule is not
	// eligible for this lookup table and has not been added.
	Add(f *rules.NetworkRule, storageIdx int64) (ok bool)

	// AppendMatching finds all matching rules from this lookup table and
	// appends them to matching.
	AppendMatching(matching []*rules.NetworkRule, r *rules.Request) (res []*rules.NetworkRule)
}
