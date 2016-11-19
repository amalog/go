package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io/ioutil"
	"os"
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

		// parse unformatted code
		term, err := readFile(base + "/" + test.Name() + "/before.ama")
		if err != nil {
			t.Errorf("can't read %s/before.ama: %s", test.Name(), err)
			continue
		}

		// fetch formatted code
		data, err := ioutil.ReadFile(base + "/" + test.Name() + "/after.ama")
		if err != nil {
			t.Errorf("can't read %s/after.ama: %s", test.Name(), err)
			continue
		}
		expected := string(data)

		// generate and compare formatted output
		got, err := asFileString(term)
		if err != nil {
			t.Errorf("trouble writing: %s", err)
			continue
		}
		if got != expected {
			t.Errorf("%s.ama:\ngot : %s\nwant: %s\n", test.Name(), got, expected)
		}

		// canonical source can be parsed
		term, err = ReadAll(strings.NewReader(got))
		if err != nil {
			t.Errorf("can't parse canonical source: %s\n%s", err, term)
			continue
		}

		// formatting it again gives the same result
		rewritten, err := asFileString(term)
		if err != nil {
			t.Errorf("trouble writing: %s", err)
			continue
		}
		if rewritten != expected {
			t.Errorf("canonical is not canonical\ngot : %s\nwant: %s\n", rewritten, expected)
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
		_, err = readFile(base + "/" + test.Name())
		if err != nil {
			t.Errorf("oops: %s", err)
			continue
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

		// comment on first line is the expected error message
		parts := strings.SplitN(amalog, "\n", 2)
		expected := parts[0][2:]
		amalog = parts[1]

		// parsing should give the expected message
		x, err := ReadAll(strings.NewReader(amalog))
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

func readFile(filename string) (Term, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadAll(file)
}

func asFileString(t Term) (string, error) {
	buf := new(bytes.Buffer)
	err := WriteAll(buf, t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
