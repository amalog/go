package scanner // import "github.com/amalog/go/scanner"

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := map[string]string{
		`hello.`: `[atom(hello) punct(.)]`,
		`hello,`: `[atom(hello) punct(,)]`,
		`3.141.`: `[num(3.141) punct(.)]`,
		`X`:      `[var(X)]`,
		`FooBar`: `[var(FooBar)]`,
		`}`:      `[punct(})]`,

		"123 hi":    `[num(123) atom(hi)]`,
		"123.4 bye": `[num(123.4) atom(bye)]`,
		"9_876":     `[num(9_876)]`,

		"foo ": `[atom(foo)]`, // trailing whitespace

		`foo{bar}`: `[atom(foo) punct({) atom(bar) punct(,) punct(})]`,

		`"hello world\n"`: `[string("hello world\n")]`,

		`use("amalog.org/std/io", Io),`: `[atom(use) punct(() string("amalog.org/std/io") punct(,) var(Io) punct(,) punct()) punct(,)]`,

		`main(W) {`: `[atom(main) punct(() var(W) punct(,) punct()) punct({)]`,

		`Io.printf(W, "Hello, world!\n"),`: `[var(Io) punct(.) atom(printf) punct(() var(W) punct(,) string("Hello, world!\n") punct(,) punct()) punct(,)]`,
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
		if err == io.EOF {
			return ts, nil
		}
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}

	return nil, nil
}
