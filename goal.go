// Goal describes any computation which produces solutions.
type Goal interface {
	// Next returns true for the first value if it succeeds in producing a solution.
	// On success, returns true for the second value if this goal can produce more
	// solutions.
	Next(Context) (bool, bool)

	// Cleanup instructs this goal to discard any resources it allocated (such as
	// goroutines) because the goal won't be called again.
	Cleanup()
}
