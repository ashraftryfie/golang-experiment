# Go Mistakes & Anti-Patterns Log

A living catalog of classic traps, common misconceptions, and hard-earned debugging lessons.

---

## 1. Variable Shadowing with `:=`

### The Trap
```go
// ANTI-PATTERN
var user *User
if authenticated {
    user, err := getUser() // SHADOWING: creates a new local 'user' scoped only inside 'if'
    if err != nil {
        return err
    }
}
// Here, 'user' is still nil!
```

### The Fix
```go
// IDIOMATIC
var user *User
var err error
if authenticated {
    user, err = getUser() // Re-uses outer 'user'
    if err != nil {
        return err
    }
}
```

> **Remember It Like This:**
> `:=` declares a *new* variable in the current block. If you intend to assign to an existing variable in an outer scope, use `=` instead.

---

## 2. Uninitialized Map Panic

### The Trap
```go
// ANTI-PATTERN
var counts map[string]int
counts["apples"] = 5 // PANIC: assignment to entry in nil map!
```

### The Fix
```go
// IDIOMATIC
counts := make(map[string]int)
counts["apples"] = 5 // Safe
```

> **Remember It Like This:**
> A `nil` slice can be appended to safely, but a `nil` map will panic on write. Always initialize maps with `make(map[K]V)` or a map literal before writing to them.

---

## 3. Goroutine Leak via Unbuffered Channel Block

### The Trap
```go
// ANTI-PATTERN
func fetchFirst(urls []string) string {
    ch := make(chan string) // Unbuffered!
    for _, url := range urls {
        go func(u string) {
            ch <- download(u) // Goroutines that finish second or third block forever!
        }(url)
    }
    return <-ch // Only reads the first result; remaining goroutines leak memory forever.
}
```

### The Fix
```go
// IDIOMATIC
func fetchFirst(urls []string) string {
    ch := make(chan string, len(urls)) // Buffered channel accommodates all writers
    for _, url := range urls {
        go func(u string) {
            ch <- download(u)
        }(url)
    }
    return <-ch // Safe: subsequent sends succeed without blocking
}
```

> **Remember It Like This:**
> If a goroutine sends on an unbuffered channel and no receiver is listening, that goroutine is stuck in memory for the lifetime of the program. Always ask: *"How does this goroutine exit?"*

---

## 4. Pointer to Slice Trap

### The Trap
```go
// ANTI-PATTERN
func appendItems(items *[]int) { // Unnecessary pointer to slice
    *items = append(*items, 10)
}
```

### The Fix
```go
// IDIOMATIC
func appendItems(items []int) []int {
    return append(items, 10)
}
```

> **Remember It Like This:**
> A slice is already a lightweight reference type pointing to a backing array. Passing a pointer to a slice (`*[]T`) is almost always an anti-pattern. Pass the slice directly and return the new slice if mutated.
