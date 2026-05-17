package coremidi

import (
	"slices"
	"testing"
)

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

func TestEntity(t *testing.T) {
	sources, err := AllSources()
	if err != nil {
		t.Fatalf("failed to get sources")
	}

	index := slices.IndexFunc(sources, func(s Source) bool {
		return s.Manufacturer() == "Apple Inc."
	})

	// test only if the source with the expected manufacturer is found, otherwise it may be a test environment issue such as
	// a missing entity
	if index == -1 {
		t.Skip("Source with the expected manufacturer not found, skipping the test")
	}

	entity := sources[index].Entity()
	if entity.Name() == "" {
		t.Fatal("failed to get entity")
	}
}
