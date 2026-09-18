# Exercise 01: Package API Design & Export Auditing (Tier 1 - Easy)

## 🎯 Problem Statement
Design a clean, encapsulated configuration package:
- The internal fields of `ServerConfig` (`host`, `port`, `timeoutSec`) must be **unexported** to prevent callers from mutating configuration behind the back of the validation logic.
- Provide an exported constructor: `NewServerConfig(host string, port int, timeoutSec int) (*ServerConfig, error)`:
  - Validates `host != ""` (`ErrEmptyHost`).
  - Validates `port` is between `1` and `65535` (`ErrInvalidPort`).
  - Validates `timeoutSec > 0` (`ErrInvalidTimeout`).
- Provide exported getter methods: `Host() string`, `Port() int`, and `Timeout() time.Duration`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement unexported fields and exported methods.
3. Run tests:
   ```powershell
   go test -v ./06-packages-modules/exercises/01-package-export
   ```
