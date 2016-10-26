package prolog

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := map[string]string{
		`hello.`: `[atom(hello) punct(.)]`,
		`3.141.`: `[num(3.141) punct(.)]`,
		`X`:      `[var(X)]`,
		`FooBar`: `[var(FooBar)]`,

		"main :-\n    hello.": `[atom(main) neck atom(hello) punct(.)]`,
		"123 hi":              `[num(123) atom(hi)]`,
		"123.4 bye":           `[num(123.4) atom(bye)]`,
		"9_876":               `[num(9_876)]`,
	}
	for prolog, expected := range tests {
		ts, err := tokens(prolog)
		if err != nil {
			t.Errorf("oops: %s", err)
			return
		}
		got := fmt.Sprintf("%s", ts)
		if got != expected {
			t.Errorf("got : %s\nwant: %s\n", got, expected)
		}
	}
}

func tokens(text string) ([]*Token, error) {
	ts := make([]*Token, 0)

	s := NewScanner(strings.NewReader(text))
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