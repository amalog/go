package scanner // import "github.com/amalog/go/scanner"

import (
	"fmt"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := map[string]string{
		`hello.`: `[atom(hello) punct(.) punct(,)]`,
		`hello,`: `[atom(hello) punct(,)]`,
		`3.141.`: `[num(3.141) punct(.) punct(,)]`,
		`X`:      `[var(X) punct(,)]`,
		`FooBar`: `[var(FooBar) punct(,)]`,
		`}`:      `[punct(}) punct(,)]`,

		"123 hi":    `[num(123) atom(hi) punct(,)]`,
		"123.4 bye": `[num(123.4) atom(bye) punct(,)]`,
		"9_876":     `[num(9_876) punct(,)]`,

		"foo ": `[atom(foo) punct(,)]`, // trailing whitespace

		`foo{bar}`: `[atom(foo) punct({) atom(bar) punct(,) punct(}) punct(,)]`,

		`"hello world\n"`: `[string("hello world\n") punct(,)]`,

		`use("amalog.org/std/io", Io),`: `[atom(use) punct(() string("amalog.org/std/io") punct(,) var(Io) punct(,) punct()) punct(,)]`,

		`main(W) {`: `[atom(main) punct(() var(W) punct(,) punct()) punct({)]`,

		`Io.printf(W, "Hello, world!\n"),`: `[var(Io) punct(.) atom(printf) punct(() var(W) punct(,) string("Hello, world!\n") punct(,) punct()) punct(,)]`,

		"foo\nbar": `[atom(foo) punct(,) atom(bar) punct(,)]`,
		"a\n\nb":   `[atom(a) punct(,) atom(b) punct(,)]`,
	}
	for prolog, expected := range tests {
		ts, err := tokens(prolog)
		if err != nil {
			t.Errorf("oops: %s", err)
			return
		}
		got := fmt.Sprintf("%s", ts)
		if got != expected {
			t.Errorf("\ngot : %s\nwant: %s\n", got, expected)
		}
	}
}

func tokens(text string) ([]*Token, error) {
	ts := make([]*Token, 0)

	s := New(strings.NewReader(text))
	for {
		t, err := s.Scan()
		if err != nil {
			return nil, err
		}
		if t.Class == Eof {
			return ts, nil
		}
		ts = append(ts, t)
	}

	return nil, nil
}
