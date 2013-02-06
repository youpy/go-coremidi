package coremidi

import "testing"

func TestNewSource(t *testing.T) {
	client, _ := NewClient("test")
	_, err := NewSource(client, "test")

	if err != nil {
		t.Fatalf("failed to create a source")
	}
}

func TestNumberOfSources(t *testing.T) {
	sources, err := AllSources()
	numberOfSources := len(sources)

	if err != nil {
		t.Fatalf("failed to get sources")
	}

	if numberOfSources <= 0 {
		t.Fatalf("invalid number of sources")
	}
}
