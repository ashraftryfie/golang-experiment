# Exercise 02: Ports and Adapters (Hexagonal Architecture)

## Objective
Implement Hexagonal Architecture (Ports & Adapters) by defining consumer-side interface contracts (Ports) and building thread-safe infrastructure implementations (Adapters).

## Concepts Covered
- **Primary / Secondary Ports**: In Go, ports are idiomatic interfaces defined in the domain/application layer where they are needed.
- **Adapters**: Concrete implementations (e.g. in-memory, SQL, external messaging) that satisfy the ports without the domain knowing their implementation details.
- **Thread Safety**: Adapters interacting with concurrent state must guard operations using mutexes (`sync.RWMutex`).
- **Defensive Cloning**: In-memory adapters must clone entities before returning or saving to prevent race conditions or external mutations bypassing domain methods.

## Requirements
1. **Domain (`Product`)**:
   - `NewProduct(id, sku string, initialStock int) (*Product, error)`
   - `DeductStock(quantity int) error`: Deducts stock; returns `ErrInsufficientStock` if `quantity > stock`.
   - `AddStock(quantity int) error`: Adds stock; returns `ErrInvalidQuantity` if `quantity <= 0`.

2. **Ports**:
   - `ProductRepository`:
     - `GetBySKU(ctx context.Context, sku string) (*Product, error)`
     - `Save(ctx context.Context, product *Product) error`
   - `AlertPublisher`:
     - `PublishLowStock(ctx context.Context, sku string, remaining int) error`

3. **Adapters**:
   - `InMemoryProductRepo`: Thread-safe implementation of `ProductRepository`. Returns `ErrProductNotFound` if SKU not present.
   - `MockAlertPublisher`: Thread-safe mock recorder implementing `AlertPublisher`, keeping a record of all published alerts.
