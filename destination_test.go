package coremidi

import "testing"

func TestNumberOfDestinations(t *testing.T) {
	destinations, _ := AllDestinations()
	numberOfDestinations := len(destinations)

	if numberOfDestinations <= 0 {
		t.Fatalf("invalid number of destinations")
	}
}
