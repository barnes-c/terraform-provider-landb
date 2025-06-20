// SPDX-FileCopyrightText: 2025 CERN
//
// SPDX-License-Identifier: GPL-3.0-or-later

package landb_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	landb "landb/internal/client"

	"github.com/stretchr/testify/require"
)

func TestDeviceCRUD(t *testing.T) {
	apiEndpoint := "https://landb.cern.ch/api/"
	clientID := "terraform-provider-landb"
	clientSecret := os.Getenv("LANDB_SSO_CLIENT_SECRET")
	audience := "production-microservice-landb-rest"
	require.NotEmpty(t, clientSecret, "environment variable LANDB_SSO_CLIENT_SECRET must be set")

	cli, err := landb.NewClient(apiEndpoint, clientID, clientSecret, audience)
	require.NoError(t, err)

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	last5 := timestamp[len(timestamp)-5:]

	deviceName := fmt.Sprintf("TF-TEST-DEVICE-%s", last5)
	serialNumber := fmt.Sprintf("SN%s", last5)
	inventoryNumber := fmt.Sprintf("INV%s", last5)

	device := landb.Device{
		Name:                 deviceName,
		SerialNumber:         serialNumber,
		InventoryNumber:      inventoryNumber,
		Tag:                  "TAG001",
		Description:          "Testing device creation",
		Zone:                 "ZONE1",
		DHCPResponse:         "ALWAYS",
		IPv4InDNSAndFirewall: true,
		IPv6InDNSAndFirewall: true,
		ManagerLock:          "NO_LOCK",
		Ownership:            "CERN",
		Location: landb.Location{
			Building: "31",
			Floor:    "1",
			Room:     "006",
		},
		Parent:       "test",
		Type:         "COMPUTER",
		Manufacturer: "APPLE MAC",
		Model:        "MACBOOK PRO 13",
		OperatingSystem: landb.OperatingSystem{
			Family:  "ANDROID",
			Version: "12",
		},
		Manager: landb.Contact{
			Type: "EGROUP",
			EGroup: landb.EGroup{
				Name:  "terraform-provider-landb",
				Email: "terraform-provider-landb@cern.ch",
			},
		},
		Responsible: landb.Contact{
			Type: "EGROUP",
			EGroup: landb.EGroup{
				Name:  "terraform-provider-landb",
				Email: "terraform-provider-landb@cern.ch",
			},
		},
		User: landb.Contact{
			Type: "EGROUP",
			EGroup: landb.EGroup{
				Name:  "terraform-provider-landb",
				Email: "terraform-provider-landb@cern.ch",
			},
		},
	}

	t.Logf("Creating device: %s", deviceName)
	createdDevice, err := cli.CreateDevice(device)
	require.NoError(t, err)
	require.Equal(t, device.Name, createdDevice.Name)

	t.Log("Reading device...")
	readDevice, err := cli.GetDevice(deviceName)
	require.NoError(t, err)
	require.Equal(t, createdDevice.Name, readDevice.Name)

	t.Log("Updating device...")
	readDevice.Description = "Updated via test"
	updatedDevice, err := cli.UpdateDevice(readDevice.Name, *readDevice)
	require.NoError(t, err)
	require.Equal(t, "Updated via test", updatedDevice.Description)

	defer func() {
		t.Logf("Deleting device: %s", updatedDevice.Name)
		err := cli.DeleteDevice(updatedDevice.Name, updatedDevice.Version)
		require.NoError(t, err)
	}()

	t.Log("Final read to confirm update...")
	finalDevice, err := cli.GetDevice(deviceName)
	require.NoError(t, err)
	require.Equal(t, "Updated via test", finalDevice.Description)
}
