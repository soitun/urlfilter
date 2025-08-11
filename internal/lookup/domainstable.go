package lookup

import (
	"strings"

	"github.com/AdguardTeam/golibs/syncutil"
	"github.com/AdguardTeam/urlfilter/filterlist"
	"github.com/AdguardTeam/urlfilter/rules"
)

// DomainsTable is a [Table] that uses domains from the $domain modifier to
// speed up the rules search.  Only the rules with $domain modifier are eligible
// for this lookup table.
type DomainsTable struct {
	// Storage for the network filtering rules.
	ruleStorage *filterlist.RuleStorage

	// subdomainsPool contains slices of strings to fill with subdomains.
	subdomainsPool *syncutil.Pool[[]string]

	// domainsIndex is the index of domains to rules that match them.
	domainsIndex map[string][]int64
}

// subdomainsEst is the estimate for the number of subdomains in a domain.
const subdomainsEst = 4

// NewDomainsTable creates a new instance of the DomainsTable.
func NewDomainsTable(rs *filterlist.RuleStorage) (s *DomainsTable) {
	return &DomainsTable{
		ruleStorage:    rs,
		subdomainsPool: syncutil.NewSlicePool[string](subdomainsEst),
		domainsIndex:   map[string][]int64{},
	}
}

// type check
var _ Table = (*DomainsTable)(nil)

// Add implements the [Table] interface for *DomainsTable.
func (d *DomainsTable) Add(f *rules.NetworkRule, storageIdx int64) (ok bool) {
	permittedDomains := f.GetPermittedDomains()
	if len(permittedDomains) == 0 {
		return false
	}

	for _, domain := range permittedDomains {
		rulesIndexes := d.domainsIndex[domain]
		rulesIndexes = append(rulesIndexes, storageIdx)
		d.domainsIndex[domain] = rulesIndexes
	}

	return true
}

// AppendMatching implements the [Table] interface for *DomainsTable.
func (d *DomainsTable) AppendMatching(
	matching []*rules.NetworkRule,
	r *rules.Request,
) (res []*rules.NetworkRule) {
	res = matching
	if r.SourceHostname == "" {
		return res
	}

	subdomainsPtr := d.subdomainsPool.Get()
	defer d.subdomainsPool.Put(subdomainsPtr)

	*subdomainsPtr = appendSubdomains((*subdomainsPtr)[:0], r.SourceHostname)
	if len(*subdomainsPtr) == 0 {
		return res
	}

	for _, domain := range *subdomainsPtr {
		matchingRules := d.domainsIndex[domain]
		for _, idx := range matchingRules {
			rule := d.ruleStorage.RetrieveNetworkRule(idx)
			if rule != nil && rule.Match(r) {
				res = append(res, rule)
			}
		}
	}

	return res
}

// appendSubdomains appends all subdomains of domain, starting from domain
// itself, to sub.  domain must be a valid, non-fully-qualified domain name.
// If domain is empty, appendSubdomains returns nil.
//
// NOTE:  Keep in sync with [netutil.Subdomains].
//
// TODO(a.garipov):  Add to golibs.
func appendSubdomains(sub []string, domain string) (res []string) {
	if domain == "" {
		return nil
	}

	res = append(sub, domain)

	for domain != "" {
		i := strings.IndexByte(domain, '.')
		if i < 0 {
			break
		}

		domain = domain[i+1:]
		res = append(res, domain)
	}

	return res
}
