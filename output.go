package ncd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Outputter describes and output handler.
type Outputter interface {
	Format() (*bytes.Buffer, error)
}

// JSONOutput will format the docket data as JSON.
type JSONOutput struct {
	DocketMapping      // docket data.
	Pretty        bool // print pretty JSON or not.
}

// NewJSONOutut returns a JSONOutput with non-pretty formatting.
func NewJSONOutput(dm DocketMapping) JSONOutput {
	return JSONOutput{
		DocketMapping: dm,
		Pretty:        false,
	}
}

// NewPrettyJSONOutut returns a JSONOutputter with pretty formatting.
func NewPrettyJSONOutput(dm DocketMapping) JSONOutput {
	return JSONOutput{
		DocketMapping: dm,
		Pretty:        true,
	}
}

// Implements Format for Outputter.
func (jo JSONOutput) Format() (*bytes.Buffer, error) {
	var j []byte
	var err error
	if jo.Pretty {
		j, err = json.MarshalIndent(jo.DocketMapping, "", "    ")
	} else {
		j, err = json.Marshal(jo.DocketMapping)
	}
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	out.Write(j)
	return out, nil
}

// TextOutput will format the docket data as plain text.
type TextOutput struct {
	DocketMapping // docket data.
}

// NewTextOutput returns a TextOutputter.
func NewTextOutput(dm DocketMapping) TextOutput {
	return TextOutput{dm}
}

// Implements Format for Outputter.
func (to TextOutput) Format() (*bytes.Buffer, error) {
	out := &bytes.Buffer{}
	for t, cas := range to.DocketMapping {
		out.WriteString(fmt.Sprintf(">> %s\n", t))
		for ca, crgs := range cas {
			ccnt := 0
			cnt := len(crgs)
			out.WriteString(fmt.Sprintf("%s\n%s\n", ca, strings.Repeat("-", len(ca))))
			for _, crg := range crgs {
				var nl string
				if ccnt != cnt {
					nl = "\n"
				}
				if crg.HasMultiple() {
					out.WriteString(fmt.Sprintf("* %s (%d counts)%s", crg.Description, crg.Count, nl))
				} else {
					out.WriteString(fmt.Sprintf("* %s%s", crg.Description, nl))
				}
				ccnt++
			}
			out.WriteString("\n")
		}
	}
	return out, nil
}
