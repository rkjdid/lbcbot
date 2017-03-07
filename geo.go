package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getPublicIP() (ip string, err error) {
	resp, err := http.Get("http://ipinfo.io/ip")
	if err != nil {
		return "", err
	}
	_, err = fmt.Fscanf(resp.Body, "%s", &ip)
	return ip, err
}

type GeoIP struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Lat         float32 `json:"latitude"`
	Lon         float32 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
	AreaCode    int     `json:"area_code"`
}

func getGeoIP(addr string) (gIP *GeoIP, err error) {
	resp, err := http.Get("https://freegeoip.net/json/" + addr)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&gIP)
	return gIP, err
}

func getLatLng() (lat, lng float32) {
	ip, err := getPublicIP()
	if err != nil {
		log.Println("couldn't get public ip addr:", err)
		return
	}
	gIP, err := getGeoIP(ip)
	if err != nil {
		log.Println("couldn't get geolocalization:", err)
	}
	return gIP.Lat, gIP.Lon
}
