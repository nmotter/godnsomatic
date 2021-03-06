# godnsomatic

DNS-O-Matic client written in the Go programming language.  The client follows the 0.9 version of the DNS-O-Matic API specifiction.

## Getting Started

git clone <https://github.com/nmotter/godnsomatic.git>

## Building the project

go build -o godnsomatic main.go

## Configuring the client

From the command line type: godnsomatic
After the initial run the config.json file is created.  See the example config.json file below.  Add your dnsomatic credintials and run godnsomatic again.

```json
    {
    "DnsomaticUsername": "dns-omatic-username",
    "DnsomaticPassword": "dns-omatic-password",
    "Hostname": [
        "all.dnsomatic.com"
    ],
    "Wildcard": "NOCHG",
    "Mx": "NOCHG",
    "Backmx": "NOCHG"
    }
```

### Property Descriptions Per DNS-Omatics Site

#### hostname

Hostname you wish to update. To update all services registered with DNS-O-Matic to the new IP address, hostname may be omitted or set to all.dnsomatic.com (useful if required by client). This field is also used for services that use different names for the unique identifier of the target being updated (ex. freedns.afraid.org, TZO). DNS-O-Matic will format the update string appropriately for each supported service at distribution.

#### wildcard

Parameter enables or disables wildcards for this host. ON enables wildcard. NOCHG value will keep current wildcard settings. Any other value will disable wildcard for hosts in update. What does wildcard do & mean in this context?

#### mx

Specifies a Mail eXchanger for use with the hostname being modified. The specified MX must resolve to an IP address, or it will be ignored. Specifying an MX of NOCHG will cause the existing MX setting to be preserved in whatever state it was previously updated.

#### backmx

Requests the MX in the previous parameter to be set up as a backup MX by listing the host itself as an MX with a lower preference value. YES activates preferred MX record pointed to hostname itself, NOCHG keeps the previous value, any other value is considered as NO and deactivates the corresponding DNS record.

For additional information please see the following DNS-O-Matic documentation: <https://www.dnsomatic.com/wiki/api>

### Scheduling Recommendations

Schedule calls to the DNS-O-Matic at least 30 minutes apart, so to not over use the service.  Once a day is probably more than enough.