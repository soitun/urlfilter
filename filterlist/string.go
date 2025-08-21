package filterlist

import (
	"strings"

	"github.com/AdguardTeam/urlfilter/rules"
)

// StringConfig represents configuration for a string-based rule list.
type StringConfig struct {
	// RulesText is a string with filtering rules (one per line).
	RulesText string

	// ID is the rule list identifier.
	ID int

	// IgnoreCosmetic tells whether to ignore cosmetic rules or not.
	IgnoreCosmetic bool
}

// String is an [Interface] implementation which stores rules within a string.
type String struct {
	rulesText      string
	id             int
	ignoreCosmetic bool
}

// NewString creates a new string-based rule list with the given configuration.
func NewString(conf *StringConfig) (s *String) {
	return &String{
		rulesText:      conf.RulesText,
		id:             conf.ID,
		ignoreCosmetic: conf.IgnoreCosmetic,
	}
}

// type check
var _ Interface = (*String)(nil)

// GetID implements the [Interface] interface for *String.
func (s *String) GetID() (id int) {
	return s.id
}

// NewScanner implements the [Interface] interface for *String.
func (s *String) NewScanner() (sc *RuleScanner) {
	return NewRuleScanner(strings.NewReader(s.rulesText), s.id, s.ignoreCosmetic)
}

// RetrieveRule implements the [Interface] interface for *String.
func (s *String) RetrieveRule(ruleIdx int) (r rules.Rule, err error) {
	if ruleIdx < 0 || ruleIdx >= len(s.rulesText) {
		return nil, ErrRuleRetrieval
	}

	endOfLine := strings.IndexByte(s.rulesText[ruleIdx:], '\n')
	if endOfLine == -1 {
		endOfLine = len(s.rulesText)
	} else {
		endOfLine += ruleIdx
	}

	line := strings.TrimSpace(s.rulesText[ruleIdx:endOfLine])
	if len(line) == 0 {
		return nil, ErrRuleRetrieval
	}

	return rules.NewRule(line, s.id)
}

// Close implements the [Interface] interface for *String.
func (s *String) Close() (err error) {
	return nil
}
