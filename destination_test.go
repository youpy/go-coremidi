package coremidi

import "testing"

func TestNumberOfDestinations(t *testing.T) {
	destinations, err := AllDestinations()
	numberOfDestinations := len(destinations)

	if err != nil {
		t.Fatalf("failed to get destinations")
	}

	if numberOfDestinations <= 0 {
		t.Fatalf("invalid number of destinations")
	}
}
