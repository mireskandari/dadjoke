package pdf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// testMode disables Wails runtime calls during tests
var testMode = false

// SetTestMode enables or disables test mode (disables Wails runtime calls)
func SetTestMode(enabled bool) {
	testMode = enabled
}

// safeEmit emits an event only if not in test mode
// This allows tests to run without requiring a Wails context
func safeEmit(ctx context.Context, eventName string, data ...interface{}) {
	if testMode {
		return
	}
	runtime.EventsEmit(ctx, eventName, data...)
}

// FormatFileSize converts bytes to human-readable format
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GenerateID creates a unique identifier
func GenerateID() string {
	return uuid.New().String()[:8]
}

// GetFileInfo returns information about a file
func GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access file: %w", err)
	}

	return &FileInfo{
		Path:     path,
		Name:     info.Name(),
		Size:     info.Size(),
		SizeText: FormatFileSize(info.Size()),
	}, nil
}

// GetPDFInfo returns detailed information about a PDF file
func GetPDFInfo(path string) (*PDFDocument, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access file: %w", err)
	}

	// Validate it's a PDF
	if !strings.HasSuffix(strings.ToLower(path), ".pdf") {
		return nil, fmt.Errorf("file is not a PDF")
	}

	// Get page count using pdfcpu
	pageCount, err := api.PageCountFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read PDF: %w", err)
	}

	return &PDFDocument{
		ID:        GenerateID(),
		Path:      path,
		Name:      info.Name(),
		PageCount: pageCount,
		Size:      info.Size(),
		SizeText:  FormatFileSize(info.Size()),
	}, nil
}

// ValidatePDF checks if a file is a valid PDF
func ValidatePDF(path string) error {
	if !strings.HasSuffix(strings.ToLower(path), ".pdf") {
		return fmt.Errorf("file is not a PDF")
	}

	if err := api.ValidateFile(path, nil); err != nil {
		// Check if it's a password-protected PDF
		if strings.Contains(err.Error(), "encrypted") || strings.Contains(err.Error(), "password") {
			return fmt.Errorf("PDF is password-protected")
		}
		return fmt.Errorf("invalid PDF: %w", err)
	}

	return nil
}

// CreateTempFile creates a temporary file with the given prefix and extension
func CreateTempFile(prefix, ext string) (string, error) {
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%s_%s%s", prefix, GenerateID(), ext))
	return tmpFile, nil
}

// CleanupTempFiles removes temporary files created during processing
func CleanupTempFiles(paths ...string) {
	for _, path := range paths {
		if strings.HasPrefix(path, os.TempDir()) {
			os.Remove(path)
		}
	}
}
