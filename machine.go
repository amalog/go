package amalog // import "github.com/amalog/go"

import (
	"os"
	"path/filepath"

	"github.com/amalog/go/term"
)

type Machine struct {
	root term.Term // term which started it all
}

func NewMachine() *Machine {
	return &Machine{}
}

// LoadRoot populates the machine's root term with the contents of a file
// and loads all of that term's module dependencies.
func (m *Machine) LoadRoot(filename string) error {
	// read root term from file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	m.root, err = term.ReadAll(file)
	if err != nil {
		return err
	}

	// find directories to search for modules
	modulePath, err := m.modulePath(filename)
	if err != nil {
		return err
	}
	_ = modulePath

	return nil
}

func (m *Machine) Call(name string, args ...term.Term) *Result {
	return nil
}

// modulePath returns a list of directories within which we should search for
// modules. The directories returned are from highest to lowest priority.  All
// directories are guaranteed to exist in the filesystem at the time this method
// was invoked.
//
// The argument should be the path to source code for the root term.
func (m *Machine) modulePath(src string) ([]string, error) {
	// start search at directories holding source code and executable
	starts := make([]string, 0, 2)
	for _, p := range []string{src, os.Args[0]} {
		p, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}
		starts = append(starts, filepath.Dir(p))
	}

	// walk back to the root looking for module directories
	paths := make([]string, 0)
	tried := make(map[string]bool)
	for _, p := range starts {
		for {
			tried[p] = true
			candidate := filepath.Join(p, "amalog_modules")
			if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
				paths = append(paths, candidate)
			} else if err != nil && !os.IsNotExist(err) {
				return nil, err // unexpected error
			}

			p = filepath.Dir(p)
			if tried[p] {
				break
			}
		}
	}

	return paths, nil
}
