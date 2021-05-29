# Request Validation for chi

specを入れるとapi.GetSwaggerメソッドも生成されるようになった。


## Schema


## Generated Code


## Run-time behavior

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
resultss
```
❯ curl -d "{\"test_multiple_of\":512}" -H "Content-Type: application/json" http://localhost:8000/posts
success

❯ curl -d "{\"test_multiple_of\":513}" -H "Content-Type: application/json" http://localhost:8000/posts
request body has an error: doesn't match the schema: Error at "/test_multiple_of": Doesn't match schema "multipleOf"
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
