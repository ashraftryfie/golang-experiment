# Exercise 03: Resolving Cyclic Import Dependencies (Tier 3 - Hard)

## 🎯 Problem Statement
In naive backend designs, the `OrderService` needs to verify if a user exists, while the `UserService` needs to fetch user orders. In many languages, they import each other. In Go, this produces a hard compiler error: `import cycle not allowed`.

### Architectural Decoupling Strategy
1. **Leaf Package `models`**: Holds clean domain entities without business logic.
2. **Interface Inversion in `orders`**: `OrderService` defines a consumer interface `UserChecker`:
   ```go
   type UserChecker interface {
       UserExists(userID string) bool
   }
   ```
3. **Dependency Injection**: Pass the `UserService` (which implements `UserChecker`) into `NewOrderService(checker UserChecker)`.
4. Result: `orders` does not import `user`, and `user` does not import `orders`! The cycle is broken into a clean Directed Acyclic Graph (DAG).

---

## 🛠️ Instructions
1. Review domain entities in [`models/models.go`](./models/models.go).
2. Wire services in [`starter_test.go`](./starter_test.go).
3. Run tests:
   ```powershell
   go test -v ./06-packages-modules/exercises/03-cyclic-decoupling/...
   ```
