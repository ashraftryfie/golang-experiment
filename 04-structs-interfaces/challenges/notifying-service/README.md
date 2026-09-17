# Challenge: Multi-Channel Notification Engine

## 🎯 Goal
Build an extensible notification dispatch engine using Go interfaces, struct embedding, and composite providers.

---

## 📋 Architecture & Types

```go
type Message struct {
    Recipient string
    Subject   string
    Body      string
}

type Notifier interface {
    Send(msg Message) error
}
```

### Providers to Implement
1. `EmailNotifier`:
   - Validates recipient has `@`.
   - Sends email notification.
2. `SMSNotifier`:
   - Validates recipient phone has at least 7 digits.
   - Sends SMS notification.
3. `MultiNotifier`:
   - Composite provider holding a slice of `Notifier`.
   - Dispatches `Send` across all registered providers, collecting any errors.

---

## 🧪 Testing
```powershell
go test -v ./04-structs-interfaces/challenges/notifying-service
```
