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

// JSONOutputter will format the docket data as JSON.
type JSONOutputter struct {
	Data   DocketMapping // docket data.
	Pretty bool          // print pretty JSON or not.
}

// Implements Format for Outputter.
func (jo JSONOutputter) Format() (*bytes.Buffer, error) {
	var j []byte
	var err error
	if jo.Pretty {
		j, err = json.MarshalIndent(jo.Data, "", "    ")
	} else {
		j, err = json.Marshal(jo.Data)
	}
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	out.Write(j)
	return out, nil
}

// TextOutputter will format the docket data as plain text.
type TextOutputter struct {
	Data DocketMapping // docket data.
}

// Implements Format for Outputter.
func (to TextOutputter) Format() (*bytes.Buffer, error) {
	out := &bytes.Buffer{}
	for t, cas := range to.Data {
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
