package main

import (
	"testing"
)

func TestGeoIP(t *testing.T) {
	pubIP, err := getPublicIP()
	if err != nil {
		t.Fatal(err)
	}
	geoIP, err := getGeoIP(pubIP)
	if err != nil {
		t.Fatal(err)
	}
	if geoIP == nil || (geoIP.Lat == 0 && geoIP.Lon == 0) {
		t.Fatal("empty geoIP")
	}
	t.Log(geoIP.Lat, geoIP.Lon)
}
