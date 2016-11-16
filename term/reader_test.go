package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	const base = "../tests/format"

	tests, err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
	}
	for _, test := range tests {
		// skip files
		if !test.IsDir() {
			continue
		}

		// fetch unformatted code
		data, err := ioutil.ReadFile(base + "/" + test.Name() + "/before.ama")
		if err != nil {
			t.Errorf("can't read %s/before.ama: %s", test.Name(), err)
			continue
		}
		amalog := string(data)

		// fetch formatted code
		data, err = ioutil.ReadFile(base + "/" + test.Name() + "/after.ama")
		if err != nil {
			t.Errorf("can't read %s/after.ama: %s", test.Name(), err)
			continue
		}
		expected := string(data)

		// generate and compare formatted output
		ts, err := terms(amalog)
		if err != nil {
			t.Errorf("oops: %s\n%s", err, amalog)
			return
		}
		var buf bytes.Buffer
		for _, term := range ts {
			buf.WriteString(term.String())
		}
		got := buf.String()
		if got != expected {
			t.Errorf("%s.ama:\ngot : %s\nwant: %s\n", test.Name(), got, expected)
		}

		// canonical source can be parsed
		ts, err = terms(got)
		if err != nil {
			t.Errorf("can't parse canonical source: %s\n%s", err, amalog)
			return
		}

		// formatting it again gives the same result
		buf = bytes.Buffer{}
		for _, term := range ts {
			buf.WriteString(term.String())
		}
		rewritten := buf.String()
		if rewritten != got {
			t.Errorf("canonical is not canonical\ngot : %s\nwant: %s\n", rewritten, got)
		}
	}
}

func TestValid(t *testing.T) {
	const base = "../tests/syntax/valid"

	tests, err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
	}
	for _, test := range tests {
		// fetch code
		data, err := ioutil.ReadFile(base + "/" + test.Name())
		if err != nil {
			t.Errorf("can't read %s: %s", test.Name(), err)
			continue
		}
		amalog := string(data)

		// parse text
		_, err = terms(amalog)
		if err != nil {
			t.Errorf("oops: %s\n%s", err, amalog)
			return
		}
	}
}

func TestInvalid(t *testing.T) {
	const base = "../tests/syntax/invalid"

	tests, err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
	}
	for _, test := range tests {
		// fetch code
		data, err := ioutil.ReadFile(base + "/" + test.Name())
		if err != nil {
			t.Errorf("can't read %s: %s", test.Name(), err)
			continue
		}
		amalog := string(data)
		parts := strings.SplitN(amalog, "\n", 2)
		expected := parts[0][2:]
		amalog = parts[1]

		x, err := terms(amalog)
		if err == nil {
			t.Errorf("no syntax error %s:\ngot : %s\nfrom: %s", test.Name(), x, amalog)
			continue
		}

		got := err.Error()
		if got != expected {
			t.Errorf("%s:\ngot : %s\nwant: %s\n%s", test.Name(), got, expected, amalog)
		}
	}
}

func terms(text string) ([]Term, error) {
	terms := make([]Term, 0)

	r := NewReader(strings.NewReader(text))
	for {
		t, err := r.Read()
		if err == io.EOF {
			return terms, nil
		}
		if err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}

	return nil, nil
}
