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
