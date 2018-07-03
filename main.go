package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	dnsOMaticMyipURL = "http://myip.dnsomatic.com"
	lastIPTxt        = "lastIp.txt"
)

var (
	key = "xixidkekdndlskdkekdnskdlfkdindkd"
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

	fmt.Println(config)

	os.Exit(0)

	// Check for ip change
	if change, ip := detectIpChange(); change == true {
		fmt.Println("The IP Changed:", ip)
		updateDNS(config, ip)
	} else {
		fmt.Println("No IP Change:", ip)
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

	resp, err := http.Get(dnsOMaticMyipURL)
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
func updateDNS(config *Config, myIp string) {

	// Build query parms
	// Example url: https://updates.dnsomatic.com/nic/update?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("https://updates.dnsomatic.com/nic/update")
	buffer.WriteString("?myip=" + myIp)
	for idx := range config.Hostnames {
		buffer.WriteString("&hostname=" + config.Hostnames[idx])
	}
	buffer.WriteString("&wildcard=" + config.Wildcard)
	buffer.WriteString("&mx=" + config.Mx)
	buffer.WriteString("&backmx=" + config.Backmx)

	client := http.Client{}
	req, err := http.NewRequest("GET", buffer.String(), nil)
	req.SetBasicAuth(config.DnsomaticUsername, config.DnsomaticPassword)
	req.Header.Set("User-Agent", "GoDNS-O-Matic/1.0")

	fmt.Println("Request:", req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Updated lastIp:", string(bodyText))
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
