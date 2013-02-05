package coremidi

import "testing"

func TestSources(t *testing.T) {
	devices, _ := AllDevices()
	device := devices[0]
	entities, _ := device.Entities()
	entity := entities[0]
	sources, _ := entity.Sources()

	if len(sources) <= 0 {
		t.Fatalf("invalid number of sources")
	}
}

func TestDestinations(t *testing.T) {
	devices, _ := AllDevices()
	device := devices[0]
	entities, _ := device.Entities()
	entity := entities[0]
	destinations, _ := entity.Destinations()

	if len(destinations) <= 0 {
		t.Fatalf("invalid number of destinations")
	}
}
