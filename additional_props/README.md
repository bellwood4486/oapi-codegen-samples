# Additional Properties

## Schema

In oapi-codegen, when `type` is specified in `additionalProperties`, it is reflected in the automatically generated
code.

```yaml
# ...snip...
components:
  schemas:
    # ...snip...
    NewPost:
      title: NewPost
      type: object
      additionalProperties: # <- Enable additional properties
        type: string
      properties:
        title:
          type: string
        content:
          type: string
      required:
        - title
```

## Generated Code

A map with the type of the value specified by `type` will be generated.

```go
package oapi

// ...snip...

type NewPost struct {
	Content              *string           `json:"content,omitempty"`
	Title                string            `json:"title"`
	AdditionalProperties map[string]string `json:"-"`
}
```

In addition, The four methods will be generated.

- Get
- Set
- UnmarshalJSON
- MarshalJSON

`Get` and `Set` provide access to the map.

In the `UnmarshalJSON` method, JSON is received once as a `json.RawMessage`, and deserialized separately from explicitly
specified fields and fields that are not.

```go
package oapi

// ...snip...

// Override default JSON handling for NewPost to handle AdditionalProperties
func (a *NewPost) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["content"]; found {
		err = json.Unmarshal(raw, &a.Content)
		if err != nil {
			return errors.Wrap(err, "error reading 'content'")
		}
		delete(object, "content")
	}

	if raw, found := object["title"]; found {
		err = json.Unmarshal(raw, &a.Title)
		if err != nil {
			return errors.Wrap(err, "error reading 'title'")
		}
		delete(object, "title")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}
```

## Run-time behavior

Values of different types specified in additionalProperties cannot be deserialized.
If the `type` is `string`, value of `b` will fail.

```json
{
  "title": "hello",
  "content": "hello world!",
  "a": "aa",
  "b": 1
}
```