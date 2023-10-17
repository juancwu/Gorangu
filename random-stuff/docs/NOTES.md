## Package Declaration

- The `main.go` should have the `package main` declaration
- Every go program starts with a package declaration
- Packages are go's way of organizing and reusing code.

## Importing Code

> Two ways of importing code

### 1. Uno!

Boring one line import statement

```go
import "fmt" // a string formating + stdio package
```

### 2. Dos!

Some what better import statement for multiple packages

```go
import (
    "fmt"
    "math"
)
```

### Built-in Doc Command

Fking go has a command for docs... `godoc <package> <function>`

Ex: `godoc fmt Println`


## Data Types

- Numbers (the typicals og)
    - Integers (starts with u is unsigned)
        - uint8
        - uint16
        - uint32
        - uint64
        - same but just `int<size>`
    - Float
        - float32
        - float64
        - complex64
        - complex128

- Strings
    - double quotes is normal string
    - back ticks is multiline string (superior 💪)

- Booleans (true/false)

## Declaring Variables

Basic syntax `var <name> <type> = <value>`

These are some other ways to declare variables

```go
// define without value
var x int32

// define without the var keyword
x := 1

// multiple variables at once
var (
    a = 1
    b = 2
    c = 3
)

// here is a constant in go
const huh string = "hUH O.o!?"
```
