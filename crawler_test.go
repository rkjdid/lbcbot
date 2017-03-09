package main

import (
	"testing"
)

func TestGetRadius(t *testing.T) {
	rad := getClosestRadius(0)
	if Radius[0]*1000 != rad {
		t.Errorf("expected %d got %d", Radius[0], rad)
	}
	rad = getClosestRadius(500)
	if Radius[len(Radius)-1]*1000 != rad {
		t.Errorf("expected %d got %d", Radius[len(Radius)-1], rad)
	}
	rad = getClosestRadius(100)
	if 100000 != rad {
		t.Errorf("expected %d got %d", 100000, rad)
	}
}

func TestQuery_Run(t *testing.T) {
	q := &Query{
		Search:   "ferrari",
		Lat:      43.297,
		Lng:      5.3875,
		RadiusKM: 200,
		Category: "voitures",
		PriceMin: 41,
		cfg:      &Config{},
	}

	items, err := q.Run()
	if err != nil {
		t.Error(err)
	}

	if len(items) < 1 {
		t.Error("no ferrari in a 200km radius around Marseille ? come on")
	}

	for k, item := range items {
		if k > 3 {
			break
		}
		t.Logf("found: %v", item)
	}
}
