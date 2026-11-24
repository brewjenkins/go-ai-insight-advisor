package password

import (
	"errors"
	"strings"
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	raw := strings.Repeat("a", 72)
	hash, err := Hash(raw)
	if err != nil || hash == raw {
		t.Fatalf("expected a password hash, got error %v", err)
	}
	if !Verify(hash, raw) {
		t.Fatal("expected matching password to verify")
	}
	if Verify(hash, "wrong password") {
		t.Fatal("accepted the wrong password")
	}
	if Verify(hash, raw+"b") {
		t.Fatal("accepted a password longer than bcrypt's byte limit")
	}
	if _, err := Hash(raw + "b"); !errors.Is(err, ErrInvalidLength) {
		t.Fatalf("expected length error, got %v", err)
	}
}
