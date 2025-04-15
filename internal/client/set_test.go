// Copyright (c) Christopher Barnes <christopher@barnes.biz>
// SPDX-License-Identifier: MPL-2.0

package landb_test

import (
	"fmt"
	"testing"
	"time"

	landb "landb/internal/client"

	"github.com/stretchr/testify/require"
)

func TestSetCRUD(t *testing.T) {
	apiEndpoint := "https://landb.cern.ch/api/"
	clientID := "terraform-provider-landb"
	clientSecret := "KWGE5p5LbPHY6nQRUNpx2EFJ91fYxYbd"
	audience := "production-microservice-landb-rest"

	cli, err := landb.NewClient(apiEndpoint, clientID, clientSecret, audience)
	require.NoError(t, err)

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	last5 := timestamp[len(timestamp)-5:]

	setName := fmt.Sprintf("TF-TEST-SET-%s", last5)

	set := landb.Set{
		Name:                 setName,
		Type:                 "INTERDOMAIN",
		NetworkDomain:        "IT-COMPUTING-NETWORK",
		Description:          "Terraform test set",
		ProjectURL:           "https://example.com",
		ReceiveNotifications: true,
		Responsible: landb.Contact{
			Type: "PERSON",
			Person: landb.Person{
				FirstName:  "Christopher",
				LastName:   "Barnes",
				Email:      "christopher.barnes@cern.ch",
				Username:   "chbarnes",
				Department: "IT",
				Group:      "CD",
			},
			EGroup: landb.EGroup{
				Name:  "ai-playground",
				Email: "christopher.barnes@cern.ch",
			},
			Reserved: landb.Reserved{
				FirstName: "Christopher",
				LastName:  "Barnes",
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

	defer func() {
		t.Logf("Deleting set: %s", createdSet.Name)
		err := cli.DeleteSet(createdSet.Name, createdSet.Version)
		require.NoError(t, err)
	}()

	t.Log("Updating set...")
	readSet.Description = "Updated set via test"
	updatedSet, err := cli.UpdateSet(readSet.Name, *readSet)
	require.NoError(t, err)
	require.Equal(t, "Updated set via test", updatedSet.Description)

	t.Log("Final read to confirm update...")
	finalSet, err := cli.GetSet(setName)
	require.NoError(t, err)
	require.Equal(t, "Updated set via test", finalSet.Description)
}
