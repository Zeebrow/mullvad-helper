package lister

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type MullvadJsonFile struct {
	Etag      string    `json:"etag"`
	Countries []Country `json:"countries"`
}

func (mjf *MullvadJsonFile) PrintOutputLikeCliForSomeReason() {
	for _, c := range mjf.Countries {
		fmt.Println(c.Name)
		for _, city := range c.Cities {
			fmt.Printf("\t%s (%s)\n", city.Name, city.Code)
			for _, r := range city.Relays {
				fmt.Printf("\t\t%s\n", r.Hostname)
			}
		}
	}
}
func (mjf *MullvadJsonFile) GetCountryCities(ctry string) ([]string, error) {
	for _, country := range mjf.Countries {
		if ctry == country.Name || ctry == country.Code {
			var cities []string
			for _, c := range country.Cities {
				returnNameFormat := fmt.Sprintf("%s (%s)", c.Name, c.Code)
				cities = append(cities, returnNameFormat)
			}
			return cities, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no such country '%s'", ctry))
}
func (mjf *MullvadJsonFile) GetCountries() []string {
	// *when* does this array get "made"?
	// has to be after len() is known, so must be
	// a function stack thing. (I think it's better to use
	// a var anyways...)
	// countries := make([]string, len(mjf.Countries))
	var countries []string
	for _, c := range mjf.Countries {
		returnNameFormat := fmt.Sprintf("%s (%s)", c.Name, c.Code)
		countries = append(countries, returnNameFormat)
	}
	return countries
}

type Country struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Cities []City `json:"cities"`
}

func (ctry *Country) getCities() []string {
	var cities []string
	for _, c := range ctry.Cities {
		returnNameFormat := fmt.Sprintf("%s (%s)", c.Name, c.Code)
		cities = append(cities, returnNameFormat)
	}
	return cities
}

type City struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Latitude  string  `json"latitude"`
	Longitude string  `json"longitude"`
	Relays    []Relay `json:"relays"`
}

func (city *City) GetRelays() []string {
	var relays []string
	for _, c := range city.Relays {
		returnNameFormat := fmt.Sprintf("%s (%s)", c.Hostname, c.Ipv4AddrIn)
		relays = append(relays, returnNameFormat)
	}
	return relays
}

type Relay struct {
	Hostname   string `json:"hostname"`     // "ae-dxb-001",
	Ipv4AddrIn string `json:"ipv4_addr_in"` // "45.9.249.34",
	Ipv6AddrIn string `json:"ipv6_addr_in"` // null,
	Owned      bool   `json:"owned"`        // false,
	Provider   string `json:"provider"`     // "M247",
}

const RELAYS_JSON = "/var/cache/mullvad-vpn/relays.json"

func GetMullvadRelays() MullvadJsonFile {
	relayJsonFile, err := os.Open(RELAYS_JSON)
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(relayJsonFile)
	var mj MullvadJsonFile
	dec.Decode(&mj)
	return mj
	// mj.printOutputLikeCliForSomeReason()
}
