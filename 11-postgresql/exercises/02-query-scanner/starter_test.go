package scanner

import (
	"errors"
	"testing"
)

type mockRows struct {
	data       [][]any
	index      int
	closed     bool
	returnErr  error
	scanErrIdx int
}

func (m *mockRows) Next() bool {
	if m.index < len(m.data) {
		m.index++
		return true
	}
	return false
}

func (m *mockRows) Scan(dest ...any) error {
	if m.index == m.scanErrIdx {
		return errors.New("scan error")
	}
	row := m.data[m.index-1]
	*dest[0].(*int64) = row[0].(int64)
	*dest[1].(*string) = row[1].(string)
	*dest[2].(*bool) = row[2].(bool)
	return nil
}

func (m *mockRows) Close() error {
	m.closed = true
	return nil
}

func (m *mockRows) Err() error {
	return m.returnErr
}

func TestScanUsers(t *testing.T) {
	mock := &mockRows{
		data: [][]any{
			{int64(1), "alice@example.com", true},
			{int64(2), "bob@example.com", false},
		},
	}

	users, err := ScanUsers(mock)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: ScanUsers is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mock.closed {
		t.Errorf("rows.Close() was not called in defer!")
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Email != "alice@example.com" || users[1].IsActive != false {
		t.Errorf("scanned data mismatch: %+v", users)
	}

	// Test iteration error propagation
	mockErr := &mockRows{
		data:      [][]any{{int64(1), "x", true}},
		returnErr: errors.New("network drop during streaming"),
	}
	_, err = ScanUsers(mockErr)
	if err == nil {
		t.Errorf("expected rows.Err() to be returned, got nil")
	}
}
