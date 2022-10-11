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

func (mjf *MullvadJsonFile) IsCity(i string) bool {
	for _, c := range mjf.Countries {
		for _, city := range c.Cities {
			if (i == city.Code) || (i == city.Name) {
				return true
			}
		}
	}
	return false
}
func (mjf *MullvadJsonFile) IsCountry(i string) bool {
	for _, c := range mjf.Countries {
		if (i == c.Code) || (i == c.Name) {
			return true
		}
	}
	return false
}

func (mjf *MullvadJsonFile) GetCountryFromCity(c string) (Country, error) {
	var rtnCountry Country
	for _, ctry := range mjf.Countries {
		for _, city := range ctry.Cities {
			if (c == city.Code) || (c == city.Name) {
				rtnCountry = ctry
				return rtnCountry, nil
			}
		}
	}
	return rtnCountry, errors.New(fmt.Sprintf("no such city '%s'\n", c))
}

func (mjf *MullvadJsonFile) GetCity(c string) (City, error) {
	var rtnCity City
	for _, ctry := range mjf.Countries {
		for _, city := range ctry.Cities {
			if (c == city.Code) || (c == city.Name) {
				rtnCity = city
				return rtnCity, nil
			}
		}
	}
	return rtnCity, errors.New(fmt.Sprintf("no such city '%s'\n", c))
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
func (mjf *MullvadJsonFile) GetCountryCodes() []string {
	var countries []string
	for _, c := range mjf.Countries {
		countries = append(countries, c.Code)
	}
	return countries
}
func (mjf *MullvadJsonFile) GetCountryNames() []string {
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
func (mjf *MullvadJsonFile) GetCountry(ctry string) (Country, error) {
	var country Country
	for _, c := range mjf.Countries {
		if (ctry == c.Name) || (ctry == c.Code) {
			country = c
			return country, nil
		}
	}
	return country, errors.New(fmt.Sprintf("No such country '%s'\n"))
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

func (ctry *Country) GetRelays() []Relay {
	var relays []Relay
	for _, c := range ctry.Cities {
		for _, r := range c.Relays {
			relays = append(relays, r)
		}
	}
	return relays
}

func (city *City) GetRelays() []Relay {
	var relays []Relay
	for _, r := range city.Relays {
		relays = append(relays, r)
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

func (relay *Relay) PrintableRelay() string {
	// us241-wireguard (23.226.135.50, 2607:fcd0:ccc0:1d05::c41f) - WireGuard, hosted by Quadranet (rented)
	owned := "rented"
	if relay.Owned {
		owned = "owned"
	}
	return fmt.Sprintf("%s (%s %s) - (TODO: proto), hosted by %s (%s)", relay.Hostname, relay.Ipv4AddrIn, relay.Ipv6AddrIn /*,*/, relay.Provider, owned)
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
