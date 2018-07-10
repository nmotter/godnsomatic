package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	lastIPTxt    = "lastIp.txt"
	dnsomaticURL = "http://myip.dnsomatic.com"
)

// Config - Handles application configuration
type Config struct {
	DnsomaticUsername string   `json:"dnsomatic_username"`
	DnsomaticPassword string   `json:"dnsomatic_password"`
	Hostname          []string `json:"hostnames"`
	Wildcard          string   `json:"wildcard"`
	Mx                string   `json:"mx"`
	Backmx            string   `json:"back_mx"`
}

var (
	config = Config{}
)

func main() {

	// Load configuration file
	loadConfig(&config)

	// Detect IP Change
	change, ip := discoverIpChange()
	if change {
		fmt.Printf("IP change detected - %s\n", ip)
		updateDNS(ip)
	} else {
		fmt.Printf("No change detected - %s\n", ip)
	}

}

// Load json configuration file.  If one is not found create one.
func loadConfig(config *Config) {

	// Create a config file if one doesn't exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {

		config := Config{
			DnsomaticUsername: "dnsUsername",
			DnsomaticPassword: "dnsPassword",
			Hostname:          []string{"all.dnsomatic.com"},
			Wildcard:          "NOCHG",
			Mx:                "NOCHG",
			Backmx:            "NOCHG",
		}

		jsonString, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			fmt.Println("Couldn't create json file!")
			fmt.Println(err)
			os.Exit(1)
		}

		err = ioutil.WriteFile("config.json", jsonString, 0755)
		if err == nil {
			fmt.Println("Missing configuration file.  Creating config.json for you now.")
			fmt.Println("Please update the configuration file and run the program again.")
			os.Exit(0)
		} else {
			fmt.Printf("Error creating configuration file %v\n", err)
			os.Exit(1)
		}
	}

	// Load config
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config.json file.  Aborting!")
		os.Exit(1)
	}
	if err := json.Unmarshal(configFile, &config); err != nil {
		fmt.Println("Error parsing config file!")
		os.Exit(1)
	}

}

// Discover IP Change
func discoverIpChange() (bool, string) {
	change := false

	// Call DNS-O-Matic my ip service
	resp, err := http.Get(dnsomaticURL)
	if err != nil {
		panic(err)
	}

	// Read body of response
	myIp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	//fmt.Println(string(myIp))

	// Create last IP file if it does not exist
	if _, err := os.Stat(lastIPTxt); os.IsNotExist(err) {
		err := ioutil.WriteFile(lastIPTxt, myIp, 0644)
		if err != nil {
			panic(err)
		}

	} else {

		// Load the last ip
		storedIp, err := ioutil.ReadFile(lastIPTxt)
		if err != nil {
			panic(err)
		}

		// Compare myIp with storedIp and detect change
		// Update lastIp.txt
		if string(myIp) != string(storedIp) {
			err := ioutil.WriteFile(lastIPTxt, myIp, 0644)
			if err != err {
				panic(err)
			}
			change = true
		}
	}

	return change, string(myIp)
}

// Call DNS-O-Matic Update Service
func updateDNS(myIp string) {

	// Build query parms
	// Example url: https://updates.dnsomatic.com/nic/update?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("https://updates.dnsomatic.com/nic/update")
	buffer.WriteString("?myip=" + myIp)
	for idx := range config.Hostname {
		buffer.WriteString("&hostname=" + config.Hostname[idx])
	}
	buffer.WriteString("&wildcard=" + config.Wildcard)
	buffer.WriteString("&mx=" + config.Mx)
	buffer.WriteString("&backmx=" + config.Backmx)

	client := http.Client{}
	req, err := http.NewRequest("GET", buffer.String(), nil)
	req.SetBasicAuth(config.DnsomaticUsername, config.DnsomaticPassword)
	req.Header.Set("User-Agent", "GoDNS-O-Matic/1.0")

	//fmt.Println("Request:", req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	fmt.Println("DNSOMATIC Response: ", string(bodyText))
}
