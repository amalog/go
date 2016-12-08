package amalog // import "github.com/amalog/go"

// Context describes the context within which a search for solutions occurs.
type Context interface {
	Trail() Trail
}
