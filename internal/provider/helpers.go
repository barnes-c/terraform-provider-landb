// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package provider

import (
	landb "landb/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func expandContact(m contactModel) landb.Contact {
	if m.Type.IsNull() || m.Type.IsUnknown() {
		return landb.Contact{}
	}

	contact := landb.Contact{
		Type: m.Type.ValueString(),
	}

	switch contact.Type {
	case "PERSON":
		if m.Person != nil {
			contact.Person = landb.Person{
				FirstName:  m.Person.FirstName.ValueString(),
				LastName:   m.Person.LastName.ValueString(),
				Email:      m.Person.Email.ValueString(),
				Username:   m.Person.Username.ValueString(),
				Department: m.Person.Department.ValueString(),
				Group:      m.Person.Group.ValueString(),
			}
		}
	case "EGROUP":
		if m.EGroup != nil {
			contact.EGroup = landb.EGroup{
				Name:  m.EGroup.Name.ValueString(),
				Email: m.EGroup.Email.ValueString(),
			}
		}
	case "RESERVED":
		if m.Reserved != nil {
			contact.Reserved = landb.Reserved{
				FirstName: m.Reserved.FirstName.ValueString(),
				LastName:  m.Reserved.LastName.ValueString(),
			}
		}
	}

	return contact
}

func flattenContact(c landb.Contact) contactModel {
	model := contactModel{
		Type: types.StringValue(c.Type),
	}

	switch c.Type {
	case "PERSON":
		model.Person = &personModel{
			FirstName:  types.StringValue(c.Person.FirstName),
			LastName:   types.StringValue(c.Person.LastName),
			Email:      types.StringValue(c.Person.Email),
			Username:   types.StringValue(c.Person.Username),
			Department: types.StringValue(c.Person.Department),
			Group:      types.StringValue(c.Person.Group),
		}
	case "EGROUP":
		model.EGroup = &egroupModel{
			Name:  types.StringValue(c.EGroup.Name),
			Email: types.StringValue(c.EGroup.Email),
		}
	case "RESERVED":
		model.Reserved = &reservedModel{
			FirstName: types.StringValue(c.Reserved.FirstName),
			LastName:  types.StringValue(c.Reserved.LastName),
		}
	}

	return model
}

type contactModel struct {
	Type     types.String   `tfsdk:"type"`
	Person   *personModel   `tfsdk:"person"`
	EGroup   *egroupModel   `tfsdk:"egroup"`
	Reserved *reservedModel `tfsdk:"reserved"`
}

type personModel struct {
	FirstName  types.String `tfsdk:"first_name"`
	LastName   types.String `tfsdk:"last_name"`
	Email      types.String `tfsdk:"email"`
	Username   types.String `tfsdk:"username"`
	Department types.String `tfsdk:"department"`
	Group      types.String `tfsdk:"group"`
}

type egroupModel struct {
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

type reservedModel struct {
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
}
