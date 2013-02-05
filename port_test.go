package coremidi

import "testing"

func TestNewOutputPort(t *testing.T) {
	client, err := NewClient("test")

	if err != nil {
		t.Fatalf("failed to create a client")
	}

	_, err = NewOutputPort(client, "test")

	if err != nil {
		t.Fatalf("failed to create a port")
	}
}
