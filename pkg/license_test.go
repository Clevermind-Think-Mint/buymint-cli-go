package license

import (
	"os"
	"testing"

	"github.com/Clevermind-Think-Mint/buymint-cli-go/internal/logger"
)

func TestValidate(t *testing.T) {
	logger.LogInit(logger.PanicLevel, false)
	// Reading license
	content, err := os.ReadFile("../test/assets/license.txt")
	if err != nil {
		t.Fatal(err)
	}
	// Building new license
	license, err := New(string(content), map[string]interface{}{
		"PublicKey": "../test/assets/public.key",
	})
	if err != nil {
		t.Fatal(err)
	}
	// Validating license against correct metadata
	correctMeta := map[string]interface{}{
		"agency": "A144109",
	}
	_, err = license.Validate(correctMeta)
	if err != nil {
		t.Errorf("Expected no error, got %v (META: got %q, want %q)", err, license.Meta, correctMeta)
	}
	// Validating license against wrong metadata
	wrongMeta := map[string]interface{}{
		"agency": "A144109",
		"Alpha":  "string to test",
		"Beta":   1,
		"Gamma":  true,
	}
	_, err = license.Validate(wrongMeta)
	if err == nil {
		t.Fatal("Validation should fail with a wrong metadata!")
	}
}
