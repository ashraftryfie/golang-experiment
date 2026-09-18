# Exercise 01: Struct Tags and Selective Serialization (Tier 1 - Easy)

## 🎯 Problem Statement
Given an internal user account struct:
1. Map fields to snake_case JSON keys (`id`, `username`, `email`, `role`, `is_active`).
2. Omit `phone_number` if it is an empty string.
3. Completely hide `password_hash` and `api_secret` from all JSON serialization (`json:"-"`).
4. Ensure unexported internal fields are not serialized.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Annotate `UserAccount` with the appropriate struct tags.
3. Implement `SerializeUser(u UserAccount) ([]byte, error)`.
4. Run tests:
   ```powershell
   go test -v ./08-files-json/exercises/01-struct-tags
   ```
