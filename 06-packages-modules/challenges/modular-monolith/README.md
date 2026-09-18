# Challenge: Modular Monolith Architecture

## 🎯 Goal
Build a clean, production-styled multi-package architecture demonstrating package separation, encapsulation with `internal/`, domain layer isolation, and public utility extraction:

```text
challenges/modular-monolith/
├── internal/
│   ├── domain/     # Core business types (pure Go, zero external dependencies)
│   └── store/      # Persistence interface and in-memory implementation
├── pkg/
│   └── validator/  # Public reusable validation helper
├── main.go         # Application wiring
└── main_test.go    # Integration verification
```

---

## 🧪 Testing
```powershell
go test -v ./06-packages-modules/challenges/modular-monolith/...
```
