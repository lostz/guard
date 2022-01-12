package guard

import (
	"fmt"
	"os"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

func loadSite(path, list string) ([]*routercommon.Domain, error) {
	geositeBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var geositeList routercommon.GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		return nil, err
	}
	for _, site := range geositeList.Entry {
		if strings.EqualFold(site.CountryCode, list) {
			return site.Domain, nil
		}
	}
	return nil, fmt.Errorf("list not found in %s ,%s", path, list)
}

func loadIP(path, country string) ([]*routercommon.CIDR, error) {
	geoipBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var geoipList routercommon.GeoIPList
	if err := proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		return nil, err
	}
	for _, geoip := range geoipList.Entry {
		if strings.EqualFold(geoip.CountryCode, country) {
			return geoip.Cidr, nil
		}
	}
	return nil, fmt.Errorf("list not found in %s ,%s", path, country)
}
