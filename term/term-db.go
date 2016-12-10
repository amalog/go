package term // import "github.com/amalog/go/term"

import "io"

type Db []Term

func NewDb(args []Term) Db {
	return Db(args)
}

func (s Db) Format(w io.Writer, style Style) {
	isRoot := style.IsRoot
	style.IsRoot = false // our children are not root

	prevName := "#" // won't impose extra newlines
	for _, t := range s {
		if isRoot {
			name := Name(t)
			if prevName != "#" && prevName != name {
				io.WriteString(w, "\n")
			}
			prevName = name
		}

		style.WriteIndent(w)
		t.Format(w, style)
	}
}

// Query represents a query to be performed against a specific Db.
type Query struct {
	db    Db
	name  string
	arity int
}

// NewQuery constructs a query to find clauses whose head unifies with the
// head of t.  A single query can be Run many times and reuses the initial
// query planning effort.
func (db Db) NewQuery(t Term) *Query {
	return &Query{
		db:    db,
		name:  Name(t),
		arity: Arity(t),
	}
}

// Run executes a query against its Db.  The resulting cursor can be used to
// fetch query results.
func (q *Query) Run() *Cursor {
	return &Cursor{
		q: q,
		i: 0,
	}
}

// Cursor represents the database location at which a query resumes.
type Cursor struct {
	q *Query
	i int
}

func (c *Cursor) Next() (Term, bool) {
	for ; c.i < len(c.q.db); c.i++ {
		candidate := c.q.db[c.i]
		if Arity(candidate) == c.q.arity && Name(candidate) == c.q.name {
			c.i++
			return candidate, c.i < len(c.q.db)
		}
	}
	return nil, false
}
