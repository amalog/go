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

		/*
			`main(W) { hi(W) }`: "main(W) {    hi(W);\n}\n",
		*/
	}
	for amalog, expected := range tests {
		ts, err := terms(amalog)
		if err != nil {
			t.Errorf("oops: %s", err)
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
