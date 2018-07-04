package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	DnsomaticUsername string
	DnsomaticPassword string
	Hostname          []string
	Wildcard          string
	Mx                string
	Backmx            string
}

func main() {

	config := Config{}

	loadConfig(&config)

}

// Load json configuration file.  If one is not found create one.
func loadConfig(config *Config) {

	// Create a config file if one doesn't exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {

		config := Config{
			DnsomaticUsername: "dnsUsername",
			DnsomaticPassword: "dnsPassword",
			Hostname:          []string{"all.dnsomatic.com"},
			Wildcard:          "dnsWildCard",
			Mx:                "mx",
			Backmx:            "backmx",
		}

		jsonString, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			fmt.Println("Couldn't create json file!")
			fmt.Println(err)
			os.Exit(1)
		}

		err = ioutil.WriteFile("config.json", jsonString, 0777)
		if err != nil {
			fmt.Println("Missing configuration file.  Creating config.json for you now.")
			fmt.Println("Please update the configuration file and run the program again.")
			os.Exit(0)
		} else {
			fmt.Println("Error creating configuration file")
			os.Exit(1)
		}
	}

	// Load config
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config.json file.  Aborting!")
		os.Exit(1)
	}
	if err = json.Unmarshal(configFile, &config); err != nil {
		fmt.Println("Error parsing config file!")
		os.Exit(1)
	}

}
