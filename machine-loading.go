package amalog // import "github.com/amalog/go"
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/amalog/go/term"
)

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

	// load dependencies needed by the root term
	err = m.loadDependencies(modulePath)
	if err != nil {
		return err
	}

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

type dependency struct {
	name     string
	variable *term.Var
}

// return module dependencies for a given term
func dependencies(t term.Term) ([]*dependency, error) {
	ds := make([]*dependency, 0)

	for _, clause := range term.Clauses(t) {
		if term.Arity(clause) == 2 && term.Name(clause) == "use" {
			name, ok := term.Arg(clause, 1).(term.String)
			if !ok {
				return nil, fmt.Errorf("in use/2, expected string as first argument got: %s", term.Arg(clause, 1))
			}
			variable, ok := term.Arg(clause, 2).(*term.Var)
			if !ok {
				return nil, fmt.Errorf("in use/2, expected var as second argument got: %s", term.Arg(clause, 2))
			}
			d := &dependency{
				name:     string(name),
				variable: variable,
			}
			ds = append(ds, d)
		}
	}

	return ds, nil
}

type loadJob struct {
	*dependency
}

type loadResult struct {
	name         string
	t            term.Term
	dependencies []*dependency
	err          error
}

// loadDependencies resolves all of the root term's dependency variables into
// actual terms.  It does the same thing, recursively, for those dependencies.
func (m *Machine) loadDependencies(path []string) error {
	workerCount := runtime.NumCPU()

	// create worker goroutines
	doneCh := make(chan struct{})
	defer close(doneCh)
	jobsCh := make(chan *loadJob)
	resultsCh := make(chan *loadResult)
	for i := 0; i < workerCount; i++ {
		go loadWorker(path, doneCh, jobsCh, resultsCh)
	}
	outgoing := jobsCh // alias so we can disable the associated select clause

	// traverse dependency graph, loading modules as we go
	needToLoad, err := dependencies(m.root)
	if err != nil {
		return err
	}
	duplicates := make([]*dependency, 0)
	isLoading := make(map[string]bool)
	alreadyLoaded := make(map[string]term.Term)
	i := 0
	for i < len(needToLoad) || len(isLoading) > 0 {
		var job *loadJob
		if i < len(needToLoad) {
			d := needToLoad[i]
			if _, ok := alreadyLoaded[d.name]; ok || isLoading[d.name] {
				duplicates = append(duplicates, d)
				i++
				continue
			}
			job = &loadJob{d}
			outgoing = jobsCh // enable select clause
		} else {
			outgoing = nil // disable select clause
		}

		select {
		case outgoing <- job:
			isLoading[job.name] = true
			i++
		case result := <-resultsCh:
			if result.err != nil {
				return result.err
			}
			alreadyLoaded[result.name] = result.t
			delete(isLoading, result.name)
			needToLoad = append(needToLoad, result.dependencies...)
		}
	}

	// bind variables for all duplicates
	for _, d := range duplicates {
		t, ok := alreadyLoaded[d.name]
		if !ok {
			log.Panicf("processing duplicate that hasn't been loaded: %#v", d)
		}
		d.variable.Value = t
	}

	return nil
}

func loadWorker(
	path []string,
	doneCh <-chan struct{},
	jobsCh <-chan *loadJob,
	resultsCh chan<- *loadResult,
) {
	for {
		select {
		case <-doneCh:
			return
		case job := <-jobsCh:
			_ = job
			resultsCh <- &loadResult{
				name:         job.name,
				t:            nil,
				dependencies: []*dependency{},
			}
		}
	}
}
