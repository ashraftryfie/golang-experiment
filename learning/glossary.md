# Go Engineering Glossary

A memory-anchored dictionary of core Go concepts, runtime mechanics, and idioms.

---

## 1. Runtime & Execution

### Goroutine
A lightweight execution thread managed entirely by the Go runtime, not the operating system kernel. Starts with an initial stack size of just ~2 KB (which grows dynamically).
> **Remember It Like This:**
> A goroutine is not an OS thread. It is a tiny task running on a pool of OS threads managed by Go's runtime scheduler.

### Channel
A typed conduit through which concurrent goroutines communicate and synchronize without explicit locks.
> **Remember It Like This:**
> Do not communicate by sharing memory; instead, share memory by communicating.

### Zero Value
The default value guaranteed to every variable upon declaration before explicit assignment: `0` for numbers, `false` for booleans, `""` for strings, and `nil` for pointers, slices, maps, channels, interfaces, and functions.
> **Remember It Like This:**
> In Go, uninitialized memory is never garbage; it is predictably zeroed out. Make the zero value useful!

---

## 2. Types & Memory

### Backing Array
The contiguous block of heap or stack memory that stores the actual elements referenced by a slice. Multiple slices can point to different overlapping sub-sections of the same backing array.
> **Remember It Like This:**
> A slice does not hold data; a slice is a 3-word window (`pointer`, `len`, `cap`) peering into a backing array.

### Value Receiver vs. Pointer Receiver
A method receiver defined on a concrete type `(t Type)` receives a full copy of the struct; mutations do not affect the caller. A pointer receiver `(t *Type)` receives the memory address; mutations modify the original instance, and pointer receivers avoid copy overhead for large structs.
> **Remember It Like This:**
> If the method mutates the struct, or if the struct is large, use a pointer receiver. If in doubt, use a pointer receiver for consistency across the type's method set.

### Implicit Interface Satisfaction
In Go, a type implements an interface simply by implementing all methods declared in that interface. There is no `implements` keyword.
> **Remember It Like This:**
> If it walks like a duck and quacks like a duck, the Go compiler treats it as a duck.

---

## 3. Error Handling & Flow

### Errors as Values
In Go, errors are ordinary values implementing the built-in `error` interface (`Error() string`). They are inspected and returned explicitly using regular control flow.
> **Remember It Like This:**
> Errors are values, not exceptions. Treat them like any other return value: inspect them, handle them, or return them wrapped.

### Error Wrapping (`%w`)
The idiom of enclosing an underlying error with contextual information using `fmt.Errorf("doing something: %w", err)`. Preserves the error chain so callers can inspect root causes via `errors.Is` and `errors.As`.
> **Remember It Like This:**
> Use `%w` when you want your caller to be able to unpack and inspect the root cause. Use `%v` when you want to hide the implementation error details.

### Defer
A keyword that schedules a function call to be executed immediately before the surrounding function returns, executed in LIFO (Last-In, First-Out) order.
> **Remember It Like This:**
> Pair resource acquisition with `defer` immediately on the next line (e.g. `file, err := os.Open(...); defer file.Close()`). Arguments are evaluated immediately, but execution waits until function exit.
