package core

import (
	"testing"
)

func TestTransaction(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := NewLogger(mockDriver)

	// Create a transaction
	txID := "test-transaction-123"
	tx := logger.NewTransaction(txID)

	// Check transaction ID
	if tx.ID() != txID {
		t.Errorf("Transaction ID = %q, want %q", tx.ID(), txID)
	}

	// Test each log level with the transaction
	tests := []struct {
		name     string
		logFunc  func(msg string, attrs ...Attributes) error
		expected Level
	}{
		{"Debug", tx.Debug, Debug},
		{"Info", tx.Info, Info},
		{"Warning", tx.Warning, Warning},
		{"Error", tx.Error, Error},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := "transaction test message"
			err := test.logFunc(msg)

			// Check no error
			if err != nil {
				t.Errorf("transaction.%s() error = %v", test.name, err)
			}

			// Check log entry was created correctly
			if len(mockDriver.Logs) != i+1 {
				t.Fatalf("Expected %d logs, got %d", i+1, len(mockDriver.Logs))
			}

			log := mockDriver.Logs[i]
			if log.Level != test.expected {
				t.Errorf("Log level = %v, want %v", log.Level, test.expected)
			}
			if log.Message != msg {
				t.Errorf("Log message = %q, want %q", log.Message, msg)
			}
			if log.TransactionID != txID {
				t.Errorf("Log transaction ID = %q, want %q", log.TransactionID, txID)
			}
		})
	}
}

func TestTransactionWithAttributes(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := NewLogger(mockDriver)

	txID := "test-transaction-456"
	tx := logger.NewTransaction(txID)

	attrs := Attributes{
		"key1": "value1",
		"key2": "value2",
	}

	err := tx.Info("transaction with attributes", attrs)

	// Check no error
	if err != nil {
		t.Errorf("transaction.Info() error = %v", err)
	}

	// Check log entry was created correctly
	if len(mockDriver.Logs) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(mockDriver.Logs))
	}

	log := mockDriver.Logs[0]
	if log.Level != Info {
		t.Errorf("Log level = %v, want %v", log.Level, Info)
	}
	if log.Message != "transaction with attributes" {
		t.Errorf("Log message = %q, want %q", log.Message, "transaction with attributes")
	}
	if log.TransactionID != txID {
		t.Errorf("Log transaction ID = %q, want %q", log.TransactionID, txID)
	}

	// Check attributes
	for k, v := range attrs {
		if log.Attrs[k] != v {
			t.Errorf("Log attribute %q = %q, want %q", k, log.Attrs[k], v)
		}
	}
}

func TestMultipleTransactions(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := NewLogger(mockDriver)

	// Create two transactions
	tx1 := logger.NewTransaction("tx-1")
	tx2 := logger.NewTransaction("tx-2")

	// Log with each transaction
	tx1.Info("tx1 message")
	tx2.Info("tx2 message")
	tx1.Warning("tx1 warning")

	// Check that we have 3 logs
	if len(mockDriver.Logs) != 3 {
		t.Fatalf("Expected 3 logs, got %d", len(mockDriver.Logs))
	}

	// Check the first log (tx1 info)
	log := mockDriver.Logs[0]
	if log.Level != Info || log.Message != "tx1 message" || log.TransactionID != "tx-1" {
		t.Errorf("Log 1 incorrect: %+v", log)
	}

	// Check the second log (tx2 info)
	log = mockDriver.Logs[1]
	if log.Level != Info || log.Message != "tx2 message" || log.TransactionID != "tx-2" {
		t.Errorf("Log 2 incorrect: %+v", log)
	}

	// Check the third log (tx1 warning)
	log = mockDriver.Logs[2]
	if log.Level != Warning || log.Message != "tx1 warning" || log.TransactionID != "tx-1" {
		t.Errorf("Log 3 incorrect: %+v", log)
	}
}
