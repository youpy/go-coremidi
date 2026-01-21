package coremidi

import "testing"

func TestNewClient(t *testing.T) {
	_, err := NewClient("test")

	if err != nil {
		t.Fatalf("invalid client")
	}
}

func TestNewClientWithNotification(t *testing.T) {
	client, err := NewClientWithNotification("test", func(notification Notification) {})
	if err != nil {
		t.Fatalf("invalid notify client")
	}

	if err := client.Close(); err != nil {
		t.Fatalf("failed to close notify client")
	}
}

func TestClientCloseIdempotent(t *testing.T) {
	client, err := NewClient("test")
	if err != nil {
		t.Fatalf("invalid client")
	}

	if err := client.Close(); err != nil {
		t.Fatalf("failed to close client")
	}

	if err := client.Close(); err != nil {
		t.Fatalf("failed to close client twice")
	}
}
