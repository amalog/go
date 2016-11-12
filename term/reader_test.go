package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	tests := map[string]string{
		`hello,`: "hello,\n",
		`X,`:     "X,\n",

		`hello, bye,`: "hello,\nbye,\n",

		// language does not expand \n inside strings
		`"hello world\n",`: "\"hello world\\n\",\n",

		`use("amalog.org/std/io",Io),`: "use(\"amalog.org/std/io\", Io),\n",

		`Io.say(W,"Hello, world!"),`: "Io.say(W, \"Hello, world!\"),\n",

		`main(W) { hi(W), },`: "main(W) {\n    hi(W),\n},\n",

		`foo() { bar() { baz, bye, }, },`: "foo {\n    bar {\n        baz,\n        bye,\n    },\n},\n",

		// comma inserted for last term inside db
		`foo{bar},`:      "foo {\n    bar,\n},\n",
		`foo{bar{hi,}},`: "foo {\n    bar {\n        hi,\n    },\n},\n",

		// comma inserted at newline and EOF
		`hello`:        "hello,\n",
		`X`:            "X,\n",
		`hello, bye`:   "hello,\nbye,\n",
		"hello\nbye":   "hello,\nbye,\n",
		"foo{bar}\n":   "foo {\n    bar,\n},\n",
		"a(x)\nb{y}\n": "a(x),\nb {\n    y,\n},\n",
		"a\n\nb":       "a,\nb,\n",

		`do { things, },`: "do {\n    things,\n},\n",
		`Loop.do { x, },`: "Loop.do {\n    x,\n},\n",

		// structs are different from atoms
		`stuff(),`:   "stuff(),\n",
		`stuff{},`:   "stuff(),\n",
		`X.stuff(),`: "X.stuff(),\n",
		`X.stuff{},`: "X.stuff(),\n",

		`thing(a) {},`: "thing(a),\n",
	}
	for amalog, expected := range tests {
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
			t.Errorf("\ngot : %s\nwant: %s\n", got, expected)
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

func TestInvalid(t *testing.T) {
	tests := map[string]string{
		`bah;`: `<input>:1:4 unexpected token: punct(;)`,

		`foo(`: `<input>:1:4 unexpected end of file`,
		`foo{`: `<input>:1:4 unexpected end of file`,
		`foo)`: `<input>:1:4 unexpected token: punct())`,
		`foo}`: `<input>:1:4 unexpected token: punct(})`,
	}
	for amalog, expected := range tests {
		x, err := terms(amalog)
		if err == nil {
			t.Errorf("no syntax error:\ngot : %s\nfrom: %s", x, amalog)
			continue
		}

		got := err.Error()
		if got != expected {
			t.Errorf("\ngot : %s\nwant: %s\n%s", got, expected, amalog)
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
