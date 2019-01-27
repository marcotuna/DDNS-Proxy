package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
)

// Hub ...
type Hub struct {
	API *cloudflare.API
}

// Setup initializes authentication with Cloudflare API
func Setup(apiKey string, apiEmail string) (*cloudflare.API, error) {

	// Construct a new API object
	api, err := cloudflare.New(apiKey, apiEmail)

	if err != nil {
		return nil, err
	}

	return api, nil
}

// UserDetails fetches current user data
func (Handler *Hub) UserDetails() (*cloudflare.User, error) {

	// Fetch user details on the account
	u, err := Handler.API.UserDetails()
	if err != nil {
		return nil, err
	}

	return &u, nil

}

// ZoneID fetches zone id from domain
func (Handler *Hub) ZoneID(domainName string) (string, error) {

	// Fetch the zone ID
	id, err := Handler.API.ZoneIDByName(domainName) // Assuming the domain exists in your Cloudflare account already
	if err != nil {
		return "", err
	}

	return id, nil

}

// ZoneDetails fetches zone details from id
func (Handler *Hub) ZoneDetails(zoneID string) (*cloudflare.Zone, error) {

	// Fetch the zone ID
	zoneDetails, err := Handler.API.ZoneDetails(zoneID)
	if err != nil {
		return nil, err
	}

	return &zoneDetails, nil

}

// CreateOrUpdateDNSRecord Receives DNS Record and Updates or Creates
func (Handler *Hub) CreateOrUpdateDNSRecord(zoneID string, dnsType string, dnsName string, dnsContent string) (bool, error) {

	// Find DNS Record
	dnsRecords, err := Handler.API.DNSRecords(zoneID, cloudflare.DNSRecord{Name: dnsName})

	if err != nil {
		return false, err
	}

	dnsRecordID := ""

	for _, dnsRecord := range dnsRecords {

		if len(dnsRecord.ID) > 0 {
			dnsRecordID = dnsRecord.ID
			break
		}
	}

	// DNS Record does not exist. Creating...
	if len(dnsRecordID) == 0 {

		// Create DNS Record
		_, err = Handler.API.CreateDNSRecord(
			zoneID,
			cloudflare.DNSRecord{
				Type:      dnsType,
				Name:      dnsName,
				Content:   dnsContent,
				Proxiable: false,
				Proxied:   false,
				TTL:       0,
				Locked:    false,
			},
		)

	} else if len(dnsRecordID) > 0 {

		// Update DNS Record
		err = Handler.API.UpdateDNSRecord(
			zoneID,
			dnsRecordID,
			cloudflare.DNSRecord{
				Type:      dnsType,
				Name:      dnsName,
				Content:   dnsContent,
				Proxiable: false,
				Proxied:   false,
				TTL:       0,
				Locked:    false,
			},
		)

	}

	if err != nil {
		return false, err
	}

	return true, err

}
