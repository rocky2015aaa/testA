package handlers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTaxInvoice(t *testing.T) {
	// Prepare the test file
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Construct the path to the test file
	filePath := filepath.Join(cwd, "..", "..", "testdata", "uploads", "27350AA.pdf")
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	ctx := context.Background()
	lines, _ := getTextFromfile(ctx, file.Name())
	invoiceData, err := parseTaxInvoice(lines)

	// Check response
	assert.NoError(t, err, "An error occurred while performing the operation")
	assert.NotNil(t, invoiceData, "An error occurred while performing the operation")
}

func TestParseBakertyInvoice(t *testing.T) {
	// Prepare the test file
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Construct the path to the test file
	filePath := filepath.Join(cwd, "..", "..", "testdata", "uploads", "bakerty.pdf")
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	ctx := context.Background()
	lines, _ := getTextFromfile(ctx, file.Name())
	invoiceData, err := parseBakertyInvoice(lines)

	// Check response
	assert.NoError(t, err, "An error occurred while performing the operation")
	assert.NotNil(t, invoiceData, "An error occurred while performing the operation")
}

func TestParseWinnersInvoice(t *testing.T) {
	// Prepare the test file
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Construct the path to the test file
	filePath := filepath.Join(cwd, "..", "..", "testdata", "uploads", "WinnersInvoice.pdf")
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	ctx := context.Background()
	lines, _ := getTextFromfile(ctx, file.Name())
	invoiceData, err := parseWinnersInvoice(lines)

	// Check response
	assert.NoError(t, err, "An error occurred while performing the operation")
	assert.NotNil(t, invoiceData, "An error occurred while performing the operation")
}

func TestParseTambaramInvoice(t *testing.T) {
	// Prepare the test file
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Construct the path to the test file
	filePath := filepath.Join(cwd, "..", "..", "testdata", "uploads", "TambaramInvoice.pdf")
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	ctx := context.Background()
	lines, _ := getTextFromfile(ctx, file.Name())
	invoiceData, err := parseTambaramInvoice(lines)

	// Check response
	assert.NoError(t, err, "An error occurred while performing the operation")
	assert.NotNil(t, invoiceData, "An error occurred while performing the operation")
}

func TestParseRPInvoice(t *testing.T) {
	// Prepare the test file
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Construct the path to the test file
	filePath := filepath.Join(cwd, "..", "..", "testdata", "uploads", "RPInvoice.pdf")
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	ctx := context.Background()
	lines, _ := getTextFromfile(ctx, file.Name())
	invoiceData, err := parseRPInvoice(lines)

	// Check response
	assert.NoError(t, err, "An error occurred while performing the operation")
	assert.NotNil(t, invoiceData, "An error occurred while performing the operation")
}
