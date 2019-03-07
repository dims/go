package boring

import "testing"

func TestNewGCMNonce(t *testing.T) {
	// Should return an error for non-standard nonce size.
	key := []byte("D249BF6DEC97B1EBD69BC4D6B3A3C49D")
	ci, err := NewAESCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	c := ci.(*aesCipher)
	_, err = c.NewGCM(gcmStandardNonceSize - 1)
	if err == nil {
		t.Error("expected error for non-standard nonce size, got none")
	}
	_, err = c.NewGCM(gcmStandardNonceSize)
	if err != nil {
		t.Errorf("expected no error for standard nonce size, got: %#v", err)
	}
}
