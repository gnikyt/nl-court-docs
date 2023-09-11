package ncd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Outputter describes and output handler.
type Outputter interface {
	Format() (*bytes.Buffer, error)
}

// JsonOutput will format the docket data as Json.
type JsonOutput struct {
	DocketMapping      // docket data.
	Pretty        bool // print pretty Json or not.
}

// NewJsonOutut returns a JsonOutput with non-pretty formatting.
func NewJsonOutput(dm DocketMapping) JsonOutput {
	return JsonOutput{
		DocketMapping: dm,
		Pretty:        false,
	}
}

// NewPrettyJsonOutut returns a JsonOutputter with pretty formatting.
func NewPrettyJsonOutput(dm DocketMapping) JsonOutput {
	return JsonOutput{
		DocketMapping: dm,
		Pretty:        true,
	}
}

// Implements Format for Outputter.
func (jo JsonOutput) Format() (*bytes.Buffer, error) {
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

// CsvOutput will format the docket data as CSV in the format of:
// [time],[case],[charge],[count].
type CsvOutput struct {
	DocketMapping // docket data.
}

// NewCsvOutput returns a CsvOutput.
func NewCsvOutput(dm DocketMapping) CsvOutput {
	return CsvOutput{dm}
}

// Implements Format for Outputter.
func (cv CsvOutput) Format() (*bytes.Buffer, error) {
	out := &bytes.Buffer{}
	csv := csv.NewWriter(out)
	for t, cas := range cv.DocketMapping {
		for ca, crgs := range cas {
			for _, crg := range crgs {
				csv.Write([]string{t, ca, crg.Description, strconv.FormatInt(int64(crg.Count), 10)})
			}
		}
	}
	csv.Flush()
	if err := csv.Error(); err != nil {
		return nil, err
	}
	return out, nil
}
