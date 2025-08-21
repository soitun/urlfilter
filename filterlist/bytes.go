package filterlist

import (
	"bytes"

	"github.com/AdguardTeam/urlfilter/rules"
)

// BytesConfig represents configuration for a bytes-based rule list.
type BytesConfig struct {
	// RulesText is the slice of bytes containing rules, each rule is separated
	// by a newline.  It must not be modified after calling [NewBytes].
	RulesText []byte

	// ID is the rule list identifier.
	ID int

	// IgnoreCosmetic tells whether to ignore cosmetic rules or not.
	IgnoreCosmetic bool
}

// Bytes is an [Interface] implementation which stores rules within a byte
// slice.
type Bytes struct {
	rulesText      []byte
	id             int
	ignoreCosmetic bool
}

// NewBytes creates a new bytes-based rule list with the given configuration.
func NewBytes(conf *BytesConfig) (s *Bytes) {
	return &Bytes{
		rulesText:      conf.RulesText,
		id:             conf.ID,
		ignoreCosmetic: conf.IgnoreCosmetic,
	}
}

// type check
var _ Interface = (*Bytes)(nil)

// GetID implements the [Interface] interface for *Bytes.
func (b *Bytes) GetID() (id int) {
	return b.id
}

// NewScanner implements the [Interface] interface for *Bytes.
func (b *Bytes) NewScanner() (sc *RuleScanner) {
	return NewRuleScanner(bytes.NewReader(b.rulesText), b.id, b.ignoreCosmetic)
}

// RetrieveRule implements the [Interface] interface for *Bytes.
func (b *Bytes) RetrieveRule(ruleIdx int) (r rules.Rule, err error) {
	if ruleIdx < 0 || ruleIdx >= len(b.rulesText) {
		return nil, ErrRuleRetrieval
	}

	endOfLine := bytes.IndexByte(b.rulesText[ruleIdx:], '\n')
	if endOfLine == -1 {
		endOfLine = len(b.rulesText)
	} else {
		endOfLine += ruleIdx
	}

	line := bytes.TrimSpace(b.rulesText[ruleIdx:endOfLine])
	if len(line) == 0 {
		return nil, ErrRuleRetrieval
	}

	return rules.NewRule(string(line), b.id)
}

// Close implements the [Interface] interface for *Bytes.
func (b *Bytes) Close() (err error) {
	return nil
}
