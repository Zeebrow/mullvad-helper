package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Zeebrow/mullvad_helper/lister"
)

/*
add alias:
alias /path/to/binary='mullvad relay list'

commands:
(empty)
	print the same output from the native app
--get-countries
	list only the available countries in which a relay is available
--get-cities=COUNTRY
	list the available cities in which a relay is available for the procided country
--get-relays=CITY
	list the available relays acailable for the procided city
*/

func doArgs() {
	var getCountries bool
	var getCities, getRelays string
	flag.BoolVar(&getCountries, "get-countries", false, "list only the available countries in which a relay is available")
	flag.StringVar(&getCities, "get-cities", "none", "list the available cities in which a relay is available for a country")
	flag.StringVar(&getRelays, "get-relays", "none", "list the available relays acailable for a city")
	flag.Parse()

	relaysJson := lister.GetMullvadRelays()
	if flag.NFlag() == 0 {
		relaysJson.PrintOutputLikeCliForSomeReason()
		os.Exit(0)
	} else if flag.NFlag() > 1 {
		fmt.Println("One at a time!")
		os.Exit(1)
	}
	if getCountries {
		countryList := relaysJson.GetCountryNames()
		for _, c := range countryList {
			fmt.Println(c)
		}
		os.Exit(0)
	}
	if getCities != "none" {
		cityList, err := relaysJson.GetCountryCities(getCities)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, c := range cityList {
			fmt.Println(c)
		}
		os.Exit(0)
	}
	if getRelays != "none" {
		if relaysJson.IsCountry(getRelays) {
			country, _ := relaysJson.GetCountry(getRelays)
			fmt.Printf("%s (%s) - All relays\n", country.Name, country.Code)
			country, err := relaysJson.GetCountry(getRelays)
			if err != nil {
				fmt.Println(err)
			}
			relays := country.GetRelays()
			for _, r := range relays {
				fmt.Printf("\t%s\n", r.PrintableRelay())
			}
			os.Exit(0)
		} else if relaysJson.IsCity(getRelays) {
			rtnCity, err := relaysJson.GetCity(getRelays)
			if err != nil {
				fmt.Println(err)
			}
			// nice to have: get country name when city is provided
			country, _ := relaysJson.GetCountryFromCity(getRelays)
			fmt.Printf("%s (%s) - %s (%s)\n", country.Name, country.Code, rtnCity.Name, rtnCity.Code)
			for _, r := range rtnCity.Relays {
				fmt.Printf("\t%s\n", r.PrintableRelay())
			}
			os.Exit(0)
		}
		os.Exit(1)

		fmt.Println(getCountries)
		fmt.Println(getCities)
		fmt.Println(getRelays)
	}
}

func main() {
	doArgs()

}
