package ncd

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OutputJSON will convert the docket data into JSON with option
// to pretty print.
func OutputJSON(d *Docket, pretty bool) ([]byte, error) {
	var j []byte
	var err error
	if pretty {
		j, err = json.MarshalIndent(d.Data, "", "    ")
	} else {
		j, err = json.Marshal(d.Data)
	}
	if err != nil {
		return nil, err
	}
	return j, nil
}

// OutputText will convert the docket data into plain text output.
func OutputText(d *Docket) string {
	var out []string
	for t, cas := range d.Data {
		out = append(out, fmt.Sprintf(">> %s", t))
		for ca, crgs := range cas {
			out = append(out, fmt.Sprintf("%s\n%s", ca, strings.Repeat("-", len(ca))))
			for _, crg := range crgs {
				var fmo string
				if crg.HasMultiple() {
					fmo = fmt.Sprintf("* %s (%d counts)", crg.Description, crg.Count)
				} else {
					fmo = fmt.Sprintf("* %s", crg.Description)
				}
				out = append(out, fmo)
			}
			out = append(out, "")
		}
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}
