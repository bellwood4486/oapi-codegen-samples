# Enum Generation

## Spec

```yaml
components:
  #...snip...
  schemas:
    #...snip...
    Color:
      type: string
      nullable: true
      enum:
        - black
        - white
        - red
        - green
        - blue
```
## Generated Code

Each constant has a type name as its prefix.

```go
package api

//...snip...

// Defines values for Color.
const (
	ColorBlack Color = "black"

	ColorBlue Color = "blue"

	ColorGreen Color = "green"

	ColorRed Color = "red"

	ColorWhite Color = "white"
)

// Color defines model for Color.
type Color string
```

## How to generate

```shell
$ make oapigen
```
