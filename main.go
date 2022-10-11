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
	var getCountryRelays, getCities, getRelays string
	flag.BoolVar(&getCountries, "get-countries", false, "list only the available countries in which a relay is available")
	flag.StringVar(&getCountryRelays, "get-country-relays", "none", "list the available relays available in a country")
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
		countryList := relaysJson.GetCountries()
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
		for _, c := range relaysJson.Countries {
			for _, city := range c.Cities {
				if (city.Name == getRelays) || (city.Code == getRelays) {
					for _, relay := range city.GetRelays() {
						fmt.Println(relay)
					}
					os.Exit(0)
				}
			}
		}
		os.Exit(1)
	}
	fmt.Println(getCountries)
	fmt.Println(getCities)
	fmt.Println(getRelays)
}

func main() {
	doArgs()

}
