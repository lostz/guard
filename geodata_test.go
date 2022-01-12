package guard

import (
	"log"
	"testing"
)

func TestLoadGeosite(t *testing.T) {
	_, err := loadSite("./geosite.dat", "cn")
	if err != nil {
		log.Fatalf(err.Error())
	}

}

func TestLoadIP(t *testing.T) {
	_, err := loadIP("./geoip.dat", "cn")
	if err != nil {
		log.Fatalf(err.Error())
	}

}
