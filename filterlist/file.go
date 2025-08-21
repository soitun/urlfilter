package filterlist

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AdguardTeam/urlfilter/rules"
)

// FileConfig represents configuration for a file-based rule list.
type FileConfig struct {
	// Path is the path to the file with rules.
	Path string

	// ID is the rule list identifier.
	ID int

	// IgnoreCosmetic tells whether to ignore cosmetic rules or not.
	IgnoreCosmetic bool
}

// File is an [Interface] implementation which stores rules within a file.
type File struct {
	// file is the file with rules.
	file *os.File

	// buffer that is used for reading from the file.
	buffer []byte

	// Mutex protects all the fields.
	//
	// TODO(d.kolyshev):  Make it private and investigate if mutex is needed.
	sync.Mutex

	// id is the rule list ID.
	id int

	// ignoreCosmetic tells whether to ignore cosmetic rules or not.
	ignoreCosmetic bool
}

// NewFile creates a new file-based rule list with the given configuration.
func NewFile(conf *FileConfig) (f *File, err error) {
	f = &File{
		id:             conf.ID,
		ignoreCosmetic: conf.IgnoreCosmetic,
		buffer:         make([]byte, readerBufferSize),
	}

	f.file, err = os.Open(filepath.Clean(conf.Path))
	if err != nil {
		return nil, err
	}

	return f, nil
}

// type check
var _ Interface = (*File)(nil)

// GetID returns the rule list identifier.
func (l *File) GetID() (id int) {
	return l.id
}

// NewScanner creates a new rules scanner that reads the list contents.
func (l *File) NewScanner() (sc *RuleScanner) {
	_, _ = l.file.Seek(0, io.SeekStart)

	return NewRuleScanner(l.file, l.id, l.ignoreCosmetic)
}

// RetrieveRule finds and deserializes rule by its index.  If there's no rule by
// that index or rule is invalid, it will return an error.
func (l *File) RetrieveRule(ruleIdx int) (r rules.Rule, err error) {
	l.Lock()
	defer l.Unlock()

	if ruleIdx < 0 {
		return nil, ErrRuleRetrieval
	}

	_, err = l.file.Seek(int64(ruleIdx), io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read line from the file.
	line, err := readLine(l.file, l.buffer)
	if err == io.EOF {
		err = nil
	}

	// Check if there were any errors while reading.
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil, ErrRuleRetrieval
	}

	return rules.NewRule(line, l.id)
}

// Close closes the underlying file.
func (l *File) Close() (err error) {
	return l.file.Close()
}

// readLine reads from the reader until '\n'.  r is the reader to read from.  b
// is the buffer to use (the idea is to reuse the same buffer when it's
// possible).
func readLine(r io.Reader, b []byte) (line string, err error) {
	for {
		var n int
		n, err = r.Read(b)
		if n > 0 {
			idx := bytes.IndexByte(b[:n], '\n')
			if idx == -1 {
				line += string(b[:n])
			} else {
				line += string(b[:idx])

				return line, err
			}
		} else {
			return line, err
		}
	}
}
