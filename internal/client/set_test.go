// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
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

func TestSetCRUD(t *testing.T) {
	apiEndpoint := "https://landb.cern.ch/api/"
	clientID := "terraform-provider-landb"
	clientSecret := os.Getenv("LANDB_SSO_CLIENT_SECRET")
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
		Description:          "Terraform test set",
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

	t.Log("Reading set...")
	readSet, err := cli.GetSet(setName)
	require.NoError(t, err)
	require.Equal(t, createdSet.Name, readSet.Name)

	t.Log("Updating set...")
	readSet.Description = "Updated set via test"
	updatedSet, err := cli.UpdateSet(readSet.Name, *readSet)
	require.NoError(t, err)
	require.Equal(t, "Updated set via test", updatedSet.Description)

	defer func() {
		t.Logf("Deleting set: %s", updatedSet.Name)
		err := cli.DeleteSet(updatedSet.Name, updatedSet.Version)
		require.NoError(t, err)
	}()

	t.Log("Final read to confirm update...")
	finalSet, err := cli.GetSet(setName)
	require.NoError(t, err)
	require.Equal(t, "Updated set via test", finalSet.Description)
}
