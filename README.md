# NL Court Docs

A script to parse the NL Court Docs with a supplied office ID and convert the docket data into readable JSON.

## Building

`go build -o build/docket ./...`.

## Running

`build/docket [office]`, example `build/docket 7`.

## Example Output

Will categorize by time, then person. With each person contianing a list of charges and the number of occurrences of those charges.

```json
{
    "09:30 AM": {
        "DOE, JOHN FOO": [
            [
                "Assault",
                "1"
            ]
        ],
        "DOE, JANE MARIA": [
            [
                "Operation of a conveyance while impaired",
                "1"
            ],
            [
                "Failure or refusal to comply with demand",
                "1"
            ],
            [
                "Dangerous operation of a conveyance",
                "1"
            ],
            [
                "Resisting or obstructing a Peace Officer",
                "1"
            ],
            [
                "Assaulting a peace officer/resisting arrest",
                "2"
            ]
        ],
      // ...
    }
    // ...
}
```
