package scanner_solution

import (
	"testing"
)

type mockRows struct {
	data      [][]any
	index     int
	closed    bool
	returnErr error
}

func (m *mockRows) Next() bool {
	if m.index < len(m.data) {
		m.index++
		return true
	}
	return false
}

func (m *mockRows) Scan(dest ...any) error {
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

func TestScanUsersSolution(t *testing.T) {
	mock := &mockRows{
		data: [][]any{
			{int64(1), "alice@example.com", true},
			{int64(2), "bob@example.com", false},
		},
	}

	users, err := ScanUsers(mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mock.closed {
		t.Errorf("rows.Close() not called in defer!")
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}
