package amalog // import "github.com/amalog/go"

// Goal describes any computation which produces solutions.
type Goal interface {
	// Next returns true for the first value if it succeeds in producing a solution.
	// On success, returns true for the second value if this goal can produce more
	// solutions.
	//
	// The two return values are usually named "ok" and "more".
	Next(Context) (bool, bool)
}
