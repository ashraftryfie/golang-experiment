# Go Syntax Cheat Sheet

A condensed reference for core Go declarations, types, and control flow.

---

## 1. Variable Declarations

```go
// 1. Explicit declaration with type (initialized to zero-value)
var count int          // count = 0
var message string     // message = ""
var isReady bool       // isReady = false

// 2. Short declaration with type inference (only allowed inside functions)
port := 8080
name := "gopher"

// 3. Block declaration (idiomatic for packages)
var (
    MaxRetries = 3
    Timeout    = 30 * time.Second
)

// 4. Constants (typed vs untyped)
const Pi = 3.14159           // Untyped float: adapts to float32 or float64
const TypedTimeout int = 60  // Strictly typed int
```

---

## 2. Control Flow

### If with Initializer
```go
// Scope of 'val' and 'err' is restricted to the if-else block
if val, err := calculate(); err != nil {
    return err
} else if val > 100 {
    fmt.Println("Large:", val)
}
```

### For Loops (The Only Loop in Go)
```go
// Standard 3-component loop
for i := 0; i < 10; i++ {
    // ...
}

// "While" style loop
for condition {
    // ...
}

// Infinite loop
for {
    // break or return to exit
}

// Range over slice / array
for idx, val := range items {
    fmt.Printf("Index: %d, Value: %v\n", idx, val)
}

// Range over map
for key, val := range lookupTable {
    fmt.Printf("Key: %s, Value: %d\n", key, val)
}
```

### Switch Statements
```go
// Expression switch (cases do NOT fall through by default)
switch status {
case "active":
    fmt.Println("Running")
case "pending", "paused": // Multi-match
    fmt.Println("Waiting")
default:
    fmt.Println("Unknown")
}

// Expressionless switch (clean replacement for long if-else chains)
switch {
case score >= 90:
    fmt.Println("Grade: A")
case score >= 80:
    fmt.Println("Grade: B")
default:
    fmt.Println("Grade: F")
}
```

---

## 3. Functions & Returns

```go
// Multiple named return values
func Divide(numerator, denominator float64) (result float64, err error) {
    if denominator == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return numerator / denominator, nil
}

// Variadic functions
func Sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}
```

---

## 4. Structs & Methods

```go
type Server struct {
    Host string
    Port int
}

// Value receiver (cannot mutate struct, receives copy)
func (s Server) Address() string {
    return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Pointer receiver (can mutate struct, avoids copy)
func (s *Server) SetPort(port int) {
    s.Port = port
}
```
