package parse_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/silk/parse"
)

func TestParseLine(t *testing.T) {
	is := is.New(t)

	var tests = []struct {
		Src  string
		Type parse.LineType
	}{{
		Src:  "",
		Type: parse.LineTypePlain,
	}, {
		Src:  "Normal text is just considered plain.",
		Type: parse.LineTypePlain,
	}, {
		Src:  "# Heading",
		Type: parse.LineTypeGroupHeading,
	}, {
		Src:  "## `POST /something`",
		Type: parse.LineTypeRequest,
	}, {
		Src:  "### Example request",
		Type: parse.LineTypePlain,
	}, {
		Src:  "* `Detail`: `123`",
		Type: parse.LineTypeDetail,
	}, {
		Src:  "* `?param=value`",
		Type: parse.LineTypeParam,
	}, {
		Src:  "  * `Detail`: `123`",
		Type: parse.LineTypeDetail,
	}, {
		Src:  "  * `?param=value`",
		Type: parse.LineTypeParam,
	}, {
		Src:  "* ?param=value",
		Type: parse.LineTypeParam,
	}, {
		Src:  "===",
		Type: parse.LineTypeSeparator,
	}, {
		Src:  "====",
		Type: parse.LineTypeSeparator,
	}, {
		Src:  "=====",
		Type: parse.LineTypeSeparator,
	}, {
		Src:  "---",
		Type: parse.LineTypeSeparator,
	}, {
		Src:  "----",
		Type: parse.LineTypeSeparator,
	}, {
		Src:  "-----",
		Type: parse.LineTypeSeparator,
	}}
	for i, test := range tests {
		l, err := parse.ParseLine(i, []byte(test.Src))
		is.NoErr(err)
		is.Equal(l.Type, test.Type)
		is.Equal(l.Bytes, []byte(test.Src))
		is.Equal(l.Number, i)
	}

}

func TestLineComments(t *testing.T) {
	is := is.New(t)
	l, err := parse.ParseLine(0, []byte(`* Key: "Value" // comments should be ignored`))
	is.NoErr(err)
	detail := l.Detail()
	is.OK(detail)
	is.Equal(detail.Key, "Key")
	is.Equal(detail.Value.Data, "Value")
	l, err = parse.ParseLine(0, []byte(`* Key: "Value" // comments should be ignored`))
	is.NoErr(err)
	is.Equal(string(l.Bytes), `* Key: "Value"`)
}

func TestLineParams(t *testing.T) {
	is := is.New(t)
	for i, line := range []string{
		"* ?key=value",
		"* `?key=value`",
		"* ?`key`=`value`",
	} {
		l, err := parse.ParseLine(i, []byte(line))
		is.NoErr(err)
		is.Equal(l.Type, parse.LineTypeParam)
		detail := l.Detail()
		is.OK(detail)
		is.Equal(detail.Key, "key")
		is.Equal(detail.Value.Data, "value")
	}
}

func TestLineDetail(t *testing.T) {
	is := is.New(t)
	l, err := parse.ParseLine(0, []byte(`* Key: "Value"`))
	is.NoErr(err)
	detail := l.Detail()
	is.Equal(detail.Key, "Key")
	is.Equal(detail.Value.Data, "Value")
}
