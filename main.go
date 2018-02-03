package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	DNS_O_MATIC_MYIP_URL = "http://myip.dnsomatic.com"
	LASTIP_TXT           = "lastIp.txt"
)

type Config struct {
	DnsomaticUsername string   `yaml:"dnsomatic_username"`
	DnsomaticPassword string   `yaml:"dnsomatic_password"`
	Hostnames         []string `yaml:"hostnames"`
	Wildcard          string   `yaml:"wildcard"`
	Mx                string   `yaml:"mx"`
	Backmx            string   `yaml:"backmx"`
}

func main() {
	config := &Config{}

	// Load configuration file
	loadConfig(config)
	fmt.Println(config.Hostnames[0])

	// Check for ip change
	if change, ip := detectIpChange(); change == true {
		fmt.Println("The IP Changed", ip)
		updateDNS(config, ip)
	} else {
		fmt.Println("No IP Change", ip)
	}
}

// Load configuration data
func loadConfig(config *Config) {

	// Load YAML config file
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	// Unmarshal YAML config file
	err2 := yaml.Unmarshal([]byte(yamlFile), config)
	if err2 != nil {
		panic(err2)
	}

}

// Call http://myip.dnsomatic.com/ to get clients ip address
func detectIpChange() (bool, string) {
	change := false

	// Call DNS-O-Matic my ip service
	resp, err := http.Get(DNS_O_MATIC_MYIP_URL)
	if err != nil {
		panic(err)
	}

	// Read body of response
	myIp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println(string(myIp))

	// Create last IP file if it does not exist
	if _, err := os.Stat(LASTIP_TXT); os.IsNotExist(err) {
		err := ioutil.WriteFile(LASTIP_TXT, myIp, 0644)
		if err != nil {
			panic(err)
		}

	} else {

		// Load the last ip
		storedIp, err := ioutil.ReadFile(LASTIP_TXT)
		if err != nil {
			panic(err)
		}

		// Compare myIp with storedIp and detect change
		// Update lastIp.txt
		if string(myIp) != string(storedIp) {
			_, err := ioutil.WriteFile(LASTIP_TXT, myIp, 0644)
			if err != err {
				panic(err)
			}
			change = true
		}

	}

	return change, string(myIp)
}

// Call DNS-O-Matic Update Service
func updateDNS(config *Config, myIp string) {

	// Build query parms
	// Example url: https://updates.dnsomatic.com/nic/update?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("https://updates.dnsomatic.com/nic/update")
	buffer.WriteString("?hostname=" + config.Hostnames[0])
	buffer.WriteString("&myip=" + myIp)
	buffer.WriteString("&wildcard=" + config.Wildcard)
	buffer.WriteString("&mx=" + config.Mx)
	buffer.WriteString("&backmx=" + config.Backmx)

	client := http.Client{}
	req, err := http.NewRequest("GET", buffer.String(), nil)
	req.SetBasicAuth(config.DnsomaticUsername, config.DnsomaticPassword)
	req.Header.Set("User-Agent", "GoDNS-O-Matic/1.0")

	fmt.Println("Request: ", req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	fmt.Println(s)
}
