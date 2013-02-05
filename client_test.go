package coremidi

import "testing"

func TestNewClient(t *testing.T) {
	_, err := NewClient("test")

	if err != nil {
		t.Fatalf("invalid client")
	}
}
