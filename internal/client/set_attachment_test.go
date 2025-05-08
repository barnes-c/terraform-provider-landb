// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package landb_test

import (
	"fmt"
	"testing"
	"time"

	landb "landb/internal/client"

	"github.com/stretchr/testify/require"
)

func TestSetAttachmentCRUD(t *testing.T) {
	apiEndpoint := "https://landb.cern.ch/api/"
	clientID := "terraform-provider-landb"
	clientSecret := "kcTLfkJF47t0NDd2cIvXHaCszuLZWdqm"
	audience := "production-microservice-landb-rest"
	require.NotEmpty(t, clientSecret, "environment variable LANDB_SSO_CLIENT_SECRET must be set")

	cli, err := landb.NewClient(apiEndpoint, clientID, clientSecret, audience)
	require.NoError(t, err)

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	last5 := timestamp[len(timestamp)-5:]
	setName := fmt.Sprintf("TF-TEST-SET-%s", last5)

	set := landb.Set{
		Name:                 setName,
		Type:                 "INTERDOMAIN",
		NetworkDomain:        "GPN",
		Description:          "Terraform test set for attachments",
		ProjectURL:           "https://example.com",
		ReceiveNotifications: true,
		Responsible: landb.Contact{
			Type: "EGROUP",
			EGroup: landb.EGroup{
				Name:  "terraform-provider-landb",
				Email: "terraform-provider-landb@cern.ch",
			},
		},
	}

	t.Logf("Creating set: %s", setName)
	createdSet, err := cli.CreateSet(set)
	require.NoError(t, err)
	require.Equal(t, set.Name, createdSet.Name)

	attachName := fmt.Sprintf("TF-ATTACH-%s", last5)
	initial := landb.SetAttachment{
		DeviceName:  attachName,
		IPv4:        "188.185.64.188",
		Description: "Initial attachment",
	}

	t.Logf("Creating attachment: %s on set %s", attachName, setName)
	createdAttach, err := cli.CreateSetAttachment(setName, initial)
	require.NoError(t, err)
	require.Equal(t, initial.DeviceName, createdAttach.DeviceName)
	require.Equal(t, initial.IPv4, createdAttach.IPv4)
	require.Equal(t, initial.IPv6, createdAttach.IPv6)

	t.Log("Listing attachments to verify creation...")
	list, err := cli.GetSetAttachments(setName)
	require.NoError(t, err)

	var found *landb.SetAttachment
	for _, a := range list {
		if a.DeviceName == attachName {
			found = &a
			break
		}
	}
	require.NotNil(t, found, "attachment should exist in list after creation")
	require.Equal(t, initial.Description, found.Description)

	found.Description = "Updated via test"
	t.Log("Updating attachment description...")
	updated, err := cli.UpdateSetAttachment(setName, attachName, *found)
	require.NoError(t, err)
	require.Equal(t, "Updated via test", updated.Description)

	t.Log("Listing attachments to verify update...")
	updatedList, err := cli.GetSetAttachments(setName)
	require.NoError(t, err)

	var updatedFound *landb.SetAttachment
	for _, a := range updatedList {
		if a.DeviceName == attachName {
			updatedFound = &a
			break
		}
	}
	require.NotNil(t, updatedFound, "attachment should exist in list after update")
	require.Equal(t, "Updated via test", updatedFound.Description)

	t.Logf("Deleting attachment: %s", attachName)
	err = cli.DeleteSetAttachment(setName, attachName)
	require.NoError(t, err)

	t.Log("Listing attachments to confirm deletion...")
	postDelList, err := cli.GetSetAttachments(setName)
	require.NoError(t, err)
	for _, a := range postDelList {
		require.NotEqual(t, attachName, a.DeviceName, "attachment should be deleted from list")
	}

	defer func() {
		t.Logf("Deleting set: %s", createdSet.Name)
		err := cli.DeleteSet(createdSet.Name, createdSet.Version)
		require.NoError(t, err)
	}()

}
