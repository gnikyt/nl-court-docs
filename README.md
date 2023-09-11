# NL Court Docs

A tool to parse the NL Court Docs with a supplied date and office ID.

## Building

`make build`.

## Running

`dist/docket [date] [office] [json|text]`

Example: `build/docket 2023-09-08 7 text`.

## Usage

```go
import (
  ncd "github.com/gnikyt/nl-court-docs"
)

// ...

// Grab today's docket for office #7.
d := ncd.NewDocket(time.Now().Format("2006-01-02"), "7", &http.client{})
res, err := d.Fetch()
if err != nil {
  log.Fatal(err)
}
if err := d.Parse(res); err != nil {
  log.Fatal(err)
}

// Output as JSON.
var out string
out, err = d.Output(ncd.NewPrettyJsonOutput(d.Data)) // or ncd.NewJsonOutput for non-pretty.
if err != nil {
  log.Fatal(err)
}
fmt.Print(out)

// Output as text.
out, _ = d.Output(ncd.NewTextOutput(d.Data))
fmt.Print(out)

// Output as CSV.
out, err = d.Output(ncd.NewCsvOutput(d.Data))
if err != nil {
  log.Fatal(err)
}
fmt.Print(out)
```

Example (JSON):

```json
{
  "09:30 AM": {
    "DOE, JOHN FOO": [
      {
        "Description": "Assault",
        "Count": "1"
      }
    ],
    "DOE, JANE MARIA": [
      {
        "Description": "Operation of a conveyance while impaired",
        "Count": "1"
      },
      {
        "Description": "Failure or refusal to comply with demand",
        "Count": "1"
      },
      {
        "Description": "Dangerous operation of a conveyance",
        "Count": "1"
      },
      {
        "Description": "Resisting or obstructing a Peace Officer",
        "Count": "1"
      },
      {
        "Description": "Assaulting a peace officer/resisting arrest",
        "Count": "2"
      }
    ],
    ["// ..."]
  },
  {"// ..."}
}
```

Example (text):

```text
>> 09:30AM
DOE, JOHN FOO
-------------
* Assult

DOE, JANE MARIA
---------------
* Operation of a conveyance while impaired
* Failure or refusal to comply with demand
* Dangerous operation of a conveyance
* Resisting or obstructing a Peace Officer
* Assaulting a peace officer/resisting arrest (2 counts)
```

Example (CSV):

```csv
09:30AM,"DOE, JOHN FOO",Assult,1
09:30AM,"DOE, JANE MARIA",Operation of a conveyance while impaired,1
09:30AM,"DOE, JANE MARIA",Failure or refusal to comply with demand,1
09:30AM,"DOE, JANE MARIA",Dangerous operation of a conveyance,1
09:30AM,"DOE, JANE MARIA",Resisting or obstructing a Peace Officer,1
09:30AM,"DOE, JANE MARIA",Assaulting a peace officer/resisting arrest,2
```

You can use the package to create your own output implementations such as saving to a database.

## API

Running `make docs` outputs the following:

```go
package ncd // import "github.com/gnikyt/nl-court-docs"

type Charge struct{ ... }
type CsvOutput struct{ ... }
    func NewCsvOutput(dm DocketMapping) CsvOutput
type Docket struct{ ... }
    func NewDocket(date string, office string, client *http.Client) *Docket
type DocketMapping = map[string]map[string][]Charge
type JsonOutput struct{ ... }
    func NewJsonOutput(dm DocketMapping) JsonOutput
    func NewPrettyJsonOutput(dm DocketMapping) JsonOutput
type Outputter interface{ ... }
type TextOutput struct{ ... }
    func NewTextOutput(dm DocketMapping) TextOutput
```
