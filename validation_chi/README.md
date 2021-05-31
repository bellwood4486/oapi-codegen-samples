# Request Validation for chi

specを入れるとapi.GetSwaggerメソッドも生成されるようになった。


## Schema


## Generated Code


## Run-time behavior

### type mismatch

yaml
```yaml
#...
properties:
  title:
    type: string
```
results
```
❯ curl -d "{\"title\":1}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/title": Field must be set to string or not be present
```

validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L985
```go
	if schemaType == "integer" {
		// ...snip...
	} else if schemaType != "" && schemaType != "number" {
		// ...snip...
	}
```
### multipleOf

yaml
```yaml
#...
      properties:
        #...
        test_multiple_of:
          type: integer
          multipleOf: 256
```
results
```
❯ curl -d "{\"test_multiple_of\":512}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_multiple_of\":513}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_multiple_of": Doesn't match schema "multipleOf"
```

validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1061
```go
    if v := schema.MultipleOf; v != nil {
        // "A numeric instance is valid only if division by this keyword's 
        //   value results in an integer." 
    	if bigFloat := big.NewFloat(value / *v); !bigFloat.IsInt() {
        	// ...snip...
        }
    }
```

`multipleOf` allows decimals.

yaml
```yaml
#...
properties:
  #...
  test_multiple_of_decimals:
    type: number
    multipleOf: 0.2
```
results
```
❯ curl -d "{\"title\":\"a\", \"test_multiple_of_decimals\":0.4}" -H "Content-Type: application/json" http://localhost:8000/posts
success
```

### maximum
yaml
```yaml
#...
      properties:
        #...
        test_maximum:
          type: integer
          maximum: 100
```
resultss
```
❯ curl -d "{\"test_maximum\":100}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_maximum\":101}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_maximum": number must be most 100
```

validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1041

```go
    // "maximum" 
    if v := schema.Max; v != nil && !(*v >= value) {
        //...snip... 
    }
```

### exclusiveMaximum
yaml
```yaml
#...
      properties:
        #...
        test_exclusive_maximum:
          type: integer
          maximum: 100
          exclusiveMaximum: true
```
resultss
```
❯ curl -d "{\"test_exclusive_maximum\":99}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_exclusive_maximum\":100}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_exclusive_maximum": number must be less than 100
```
validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1007
```go
    // "exclusiveMaximum"
    if v := schema.ExclusiveMax; v && !(*schema.Max > value) {
    	// ...snip...
    }
```

### minimum
yaml
```yaml
#...
      properties:
        #...
        test_minimum:
          type: integer
          minimum: 10
```
results
```
❯ curl -d "{\"test_minimum\":10}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_minimum\":9}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_minimum": number must be at least 10
```
validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1024
```go
	// "minimum"
	if v := schema.Min; v != nil && !(*v <= value) {
		// ...snip...
	}
```

### exclusive_minimum
yaml
```yaml
#...
      properties:
        #...
        test_exclusive_minimum:
          type: integer
          minimum: 10
          exclusiveMinimum: true
```
results
```
❯ curl -d "{\"test_exclusive_minimum\":11}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_exclusive_minimum\":10}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_exclusive_minimum": number must be more than 10
```
validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L990
```go
    // "exclusiveMinimum"
    if v := schema.ExclusiveMin; v && !(*schema.Min < value) {
    	// ...snip...
    }
```

### max_length
yaml
```yaml
#...
      properties:
        #...
        test_max_length:
          type: string
          maxLength: 10
```
results
```
❯ curl -d "{\"test_max_length\":\"1234567890\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_max_length\":\"1234567890A\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_max_length": maximum string length is 10
```
validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1124
```go
	// "minLength" and "maxLength"
	minLength := schema.MinLength
	maxLength := schema.MaxLength
	if minLength != 0 || maxLength != nil {
		// JSON schema string lengths are UTF-16, not UTF-8!
		length := int64(0)
		for _, r := range value {
			if utf16.IsSurrogate(r) {
				length += 2
			} else {
				length++
			}
		}
		// ...snip...
		if maxLength != nil && length > int64(*maxLength) {
			// ...snip...
		}
	}
```

`// JSON schema string lengths are UTF-16, not UTF-8!` :think_face:

### min_length
yaml
```yaml
#...
      properties:
        #...
        test_min_length:
          type: string
          minLength: 5
```
results
```
❯ curl -d "{\"test_min_length\":\"12345\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_min_length\":\"1234\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_min_length": minimum string length is 5
```
validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1109
```go
	// "minLength" and "maxLength"
	minLength := schema.MinLength
	maxLength := schema.MaxLength
	if minLength != 0 || maxLength != nil {
		// JSON schema string lengths are UTF-16, not UTF-8!
		length := int64(0)
		for _, r := range value {
			if utf16.IsSurrogate(r) {
				length += 2
			} else {
				length++
			}
		}
		if minLength != 0 && length < int64(minLength) {
			// ...snip...
		}
		// ...snip...
	}
```

### test_pattern
yaml
```yaml
#...
      properties:
        #...
        test_pattern:
          type: string
          pattern: '^[A-Za-z]+'
```
results
```
❯ curl -d "{\"test_pattern\":\"a\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_pattern\":\"0a\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_pattern": string doesn't match the regular expression "^[A-Za-z]+"
```
validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1142
```go
    // "pattern"
    if pattern := schema.Pattern; pattern != "" && schema.compiledPattern == nil {
        var err error
        if schema.compiledPattern, err = regexp.Compile(pattern); err != nil {
        	// ...snip...
        }
        if cp := schema.compiledPattern; cp != nil && !cp.MatchString(value) {
        	// ...snip...
        }
    }
```

### max_items
yaml
```yaml
#...
      properties:
        #...
        test_max_items:
          type: array
          maxItems: 5
          items:
            type: integer
```
results
```
❯ curl -d "{\"test_max_items\":[1,2,3,4,5]}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_max_items\":[1,2,3,4,5,6]}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_max_items": maximum number of items is 5
```

### min_items
yaml
```yaml
#...
      properties:
        #...
        test_min_items:
          type: array
          minItems: 2
          items:
            type: integer
```
results
```
❯ curl -d "{\"test_min_items\":[1,2]}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_min_items\":[1]}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_min_items": minimum number of items is 2
```

### unique_items
yaml
```yaml
#...
      properties:
        #...
        test_unique_items:
          type: array
          uniqueItems: true
          items:
            type: integer
```
results
```
❯ curl -d "{\"test_unique_items\":[1,2,3]}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_unique_items\":[1,2,3,2]}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_unique_items": duplicate items found
```

### enum
yaml
```yaml
#...
      properties:
        #...
        test_enum:
          type: string
          enum:
            - dog
            - cat
            - bird
```
results
```
❯ curl -d "{\"test_enum\":\"dog\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_enum\":\"pig\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_enum": value is not one of the allowed values
```

### format(string)

Defined in [OpenAPIv3.0.3](https://swagger.io/specification/#data-types)

|Format|Support|Default or Optional|
|------|-------|--------|
|byte  |o      |Default |
|binary|x      |-       |
|date  |o      |Default |
|date-time|o   |Default |
|password|x    |-       |

Defined in [JSON Schema](https://json-schema.org/draft/2020-12/json-schema-validation.html#rfc.section.7.3)

|Format         |Support|Default or Optional|
|---------------|-------|--------|
|email          |o      |Default |
|idn-email      |x      |-       |
|hostname       |x      |-       |
|idn-hostname   |x      |-       |
|ipv4           |o      |[Optional](https://github.com/getkin/kin-openapi/blob/17153345908503543b50b7b6409f9d030bae0beb/openapi3/schema_formats.go#L98)|
|ipv6           |o      |[Optional](https://github.com/getkin/kin-openapi/blob/17153345908503543b50b7b6409f9d030bae0beb/openapi3/schema_formats.go#L103)|
|uri            |x      |-       |
|uri-reference  |x      |-       |
|iri            |x      |-       |
|iri-reference  |x      |-       |
|uuid           |x      |-       |
|uri-template   |x      |-       |
|json-pointer   |x      |-       |
|relative-json-pointer|x|-       |
|regex          |x      |-       |

kin-openapi set up the default validators in [`init` function](https://github.com/getkin/kin-openapi/blob/17153345908503543b50b7b6409f9d030bae0beb/openapi3/schema_formats.go#L80).

To use the optional validators, Call a function explicitly.

validation code
https://github.com/getkin/kin-openapi/blob/93b779808793a8a6b54ffc1f87ba17d0ffa12b70/openapi3/schema.go#L1175
```go
    // "format"
    var formatErr string
    if format := schema.Format; format != "" {
        if f, ok := SchemaStringFormats[format]; ok {
            switch {
            case f.regexp != nil && f.callback == nil:
                if cp := f.regexp; !cp.MatchString(value) {
                    formatErr = fmt.Sprintf("string doesn't match the format %q (regular expression %q)", format, cp.String())
                }
            case f.regexp == nil && f.callback != nil:
                if err := f.callback(value); err != nil {
                    formatErr = err.Error()
                }
            default:
                formatErr = fmt.Sprintf("corrupted entry %q in SchemaStringFormats", format)
            }
        }
    }
    if formatErr != "" {
    	// ...ship...
    }
```

#### byte

yaml
```yaml
      properties:
        #...
        test_format_byte:
          type: string
          format: byte
```
results
```
❯ curl -d "{\"test_format_byte\":\"aGVsbG8gd29ybGQ=\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_format_byte\":\"😀\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_format_byte": string doesn't match the format "byte" (regular expression "(^$|^[a-zA-Z0-9+/\\-_]*=*$)")
```

#### date

yaml
```yaml
      properties:
        #...
        test_format_date:
          type: string
          format: date
```
results
```
❯ curl -d "{\"test_format_date\":\"2021-05-31\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_format_date\":\"😀\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_format_date": string doesn't match the format "date" (regular expression "^[0-9]{4}-(0[0-9]|10|11|12)-([0-2][0-9]|30|31)$")
```

#### date-time

yaml
```yaml
      properties:
        #...
        test_format_datetime:
          type: string
          format: date-time
```
results
```
❯ curl -d "{\"test_format_datetime\":\"2021-05-31T16:27:35+09:00\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_format_datetime\":\"😀\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_format_datetime": string doesn't match the format "date-time" (regular expression "^[0-9]{4}-(0[0-9]|10|11|12)-([0-2][0-9]|30|31)T[0-9]{2}:[0-9]{2}:[0-9]{2}(.[0-9]+)?(Z|(\\+|-)[0-9]{2}:[0-9]{2})?$")
```

#### email

yaml
```yaml
      properties:
        #...
        test_format_datetime:
          type: string
          format: email
```
results
```
❯ curl -d "{\"test_format_email\":\"test@example.com\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_format_email\":\"😀\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_format_email": string doesn't match the format "email" (regular expression "^[^@]+@[^@<>\",\\s]+$")
```

#### ipv4

Call `DefineIPv4Format` function for validation.
```go
func main() {
	// ...snip...
	openapi3.DefineIPv4Format()

	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	// ...snip...
}
```
yaml
```yaml
      properties:
        #...
        test_format_datetime:
          type: string
          format: ipv4
```
results
```
❯ curl -d "{\"test_format_ipv4\":\"127.0.0.1\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_format_ipv4\":\"😀\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_format_ipv4": Not an IP address
```

#### ipv6

Call `DefineIPv6Format` function for validation.
```go
func main() {
// ...snip...
openapi3.DefineIPv6Format()

r := chi.NewRouter()
r.Use(middleware.OapiRequestValidator(swagger))
// ...snip...
}
```

yaml
```yaml
      properties:
        #...
        test_format_datetime:
          type: string
          format: ipv6
```
results
```
❯ curl -d "{\"test_format_ipv6\":\"2001:0db8:0000:0000:1235:0000:0000:0abc\"}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_format_ipv6\":\"😀\"}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_format_ipv6": Not an IP address
```

#### custom(regex base)


#### custom(callback base)

例：uuid
