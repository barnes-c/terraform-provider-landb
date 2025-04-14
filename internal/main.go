package main

import (
	"fmt"
	"log"

	landb "landb/internal/client"
)

func main() {
	apiEndpoint := "https://landb.cern.ch/api/"
	clientID := "terraform-provider-landb"
	clientSecret := "KWGE5p5LbPHY6nQRUNpx2EFJ91fYxYbd"
	audience := "production-microservice-landb-rest"

	cli, err := landb.NewClient(apiEndpoint, clientID, clientSecret, audience)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	device := &landb.Device{
		Name:                "test-device-123",
		SerialNumber:        "SN12345",
		InventoryNumber:     "INV67890",
		Tag:                 "TAG001",
		Description:         "Testing device creation",
		Zone:                "ZONE1",
		DHCPResponse:        "ALWAYS",
		IPv4InDNSAndFirewall: true,
		IPv6InDNSAndFirewall: true,
		ManagerLock:         "NO_LOCK",
		Ownership:           "CERN",
		Location: landb.Location{
			Building: "31",
			Floor:    "1",
			Room:     "006",
		},
		Parent:       "",
		Type:         "BRIDGE",
		Manufacturer: "Cisco",
		Model:        "ABC123",
		OperatingSystem: landb.OperatingSystem{
			Family:  "ANDROID",
			Version: "12",
		},
		Manager: landb.Contact{
			Type: "PERSON",
			Person: landb.Person{
				FirstName: "Christopher",
				LastName:  "Barnes",
				Email:     "christopher.barnes@cern.ch",
				Username:  "chbarnes",
				Department: "IT",
				Group:      "CD",
			},
		},
		Responsible: landb.Contact{
			Type: "PERSON",
			Person: landb.Person{
				FirstName: "Christopher",
				LastName:  "Barnes",
				Email:     "christopher.barnes@cern.ch",
				Username:  "chbarnes",
				Department: "IT",
				Group:      "CD",
			},
		},
		User: landb.Contact{
			Type: "PERSON",
			Person: landb.Person{
				FirstName: "Christopher",
				LastName:  "Barnes",
				Email:     "christopher.barnes@cern.ch",
				Username:  "chbarnes",
				Department: "IT",
				Group:      "CD",
			},
		},
		Version: 0,
	}

	newDevice, err := cli.CreateDevice(*device)
	if err != nil {
		log.Fatalf("failed to create device: %v", err)
	}
	fmt.Printf("Device created successfully: %+v\n", newDevice)

	device, err = cli.GetDevice("test-device-123")
	if err != nil {
		log.Fatalf("failed to get device: %v", err)
	}
	fmt.Printf("Found device: %+v\n", device)

	device.Description = "Updated description"

	updatedDevice, err := cli.UpdateDevice(device.Name, *device)
	if err != nil {
		log.Fatalf("failed to update device: %v", err)
	}
	fmt.Printf("Device updated successfully: %+v\n", updatedDevice)

	device, err = cli.GetDevice("test-device-123")
	if err != nil {
		log.Fatalf("failed to get device: %v", err)
	}

	err = cli.DeleteDevice(device.Name, device.Version)
	if err != nil {
		log.Fatalf("failed to delete device: %v", err)
	}
	fmt.Println("Device deleted successfully")
}
