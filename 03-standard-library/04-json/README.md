# encoding/json

> 📝 **Planned.** No code yet. This is what I plan to add here.

Turn Go values into JSON and back. Basic struct tags are already shown in [structs/03-json-tags](../../01-fundamentals/07-structs/03-json-tags/).

## Examples to add

- `01-marshal-unmarshal`: struct to JSON and back, with `omitempty` and `"-"`
- `02-dynamic`: read JSON into `map[string]any`
- `03-streaming`: `json.Decoder` and `json.Encoder`
- `04-custom`: write your own `MarshalJSON`

## Things to remember

- **Only fields that start with a capital letter are used.**
- **Numbers read into `any` become `float64`.**
- `DisallowUnknownFields` rejects JSON with extra fields.

## Go version notes

- The `omitzero` tag needs Go 1.24, so it can't be used here yet.
