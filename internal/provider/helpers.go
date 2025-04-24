// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package provider

import (
	"context"

	landb "landb/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	personAttrTypes = map[string]attr.Type{
		"first_name": types.StringType,
		"last_name":  types.StringType,
		"email":      types.StringType,
		"username":   types.StringType,
		"department": types.StringType,
		"group":      types.StringType,
	}
	egroupAttrTypes = map[string]attr.Type{
		"name":  types.StringType,
		"email": types.StringType,
	}
	reservedAttrTypes = map[string]attr.Type{
		"first_name": types.StringType,
		"last_name":  types.StringType,
	}
	contactAttrTypes = map[string]attr.Type{
		"type":     types.StringType,
		"person":   types.ObjectType{AttrTypes: personAttrTypes},
		"egroup":   types.ObjectType{AttrTypes: egroupAttrTypes},
		"reserved": types.ObjectType{AttrTypes: reservedAttrTypes},
	}
)

func contactSchemaBlock(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "One of PERSON, EGROUP, or RESERVED",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"person": schema.SingleNestedAttribute{
				Description: "Details if type == PERSON",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"first_name": schema.StringAttribute{Optional: true, Computed: true},
					"last_name":  schema.StringAttribute{Optional: true, Computed: true},
					"email":      schema.StringAttribute{Optional: true, Computed: true},
					"username":   schema.StringAttribute{Optional: true, Computed: true},
					"department": schema.StringAttribute{Optional: true, Computed: true},
					"group":      schema.StringAttribute{Optional: true, Computed: true},
				},
			},
			"egroup": schema.SingleNestedAttribute{
				Description: "Details if type == EGROUP",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"name":  schema.StringAttribute{Optional: true, Computed: true},
					"email": schema.StringAttribute{Optional: true, Computed: true},
				},
			},
			"reserved": schema.SingleNestedAttribute{
				Description: "Details if type == RESERVED",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"first_name": schema.StringAttribute{Optional: true, Computed: true},
					"last_name":  schema.StringAttribute{Optional: true, Computed: true},
				},
			},
		},
	}
}

func expandContactObject(ctx context.Context, o types.Object) (landb.Contact, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o.IsNull() || o.IsUnknown() {
		return landb.Contact{}, diags
	}

	elems := o.Attributes()
	typRaw, ok := elems["type"]
	if !ok {
		diags.AddError("Bad contact block", "missing \"type\" field")
		tflog.Error(ctx, "expandContactObject: missing type", map[string]interface{}{})
		return landb.Contact{}, diags
	}
	typ := typRaw.(types.String).ValueString()

	contact := landb.Contact{Type: typ}

	switch typ {
	case "PERSON":
		pRaw := elems["person"].(types.Object)
		if !pRaw.IsNull() && !pRaw.IsUnknown() {
			p := pRaw.Attributes()
			contact.Person = landb.Person{
				FirstName:  p["first_name"].(types.String).ValueString(),
				LastName:   p["last_name"].(types.String).ValueString(),
				Email:      p["email"].(types.String).ValueString(),
				Username:   p["username"].(types.String).ValueString(),
				Department: p["department"].(types.String).ValueString(),
				Group:      p["group"].(types.String).ValueString(),
			}
		}
	case "EGROUP":
		eRaw := elems["egroup"].(types.Object)
		if !eRaw.IsNull() && !eRaw.IsUnknown() {
			e := eRaw.Attributes()
			contact.EGroup = landb.EGroup{
				Name:  e["name"].(types.String).ValueString(),
				Email: e["email"].(types.String).ValueString(),
			}
		}
	case "RESERVED":
		rRaw := elems["reserved"].(types.Object)
		if !rRaw.IsNull() && !rRaw.IsUnknown() {
			r := rRaw.Attributes()
			contact.Reserved = landb.Reserved{
				FirstName: r["first_name"].(types.String).ValueString(),
				LastName:  r["last_name"].(types.String).ValueString(),
			}
		}
	}

	return contact, diags
}

func flattenContactObject(c landb.Contact) types.Object {
	elems := map[string]attr.Value{
		"type": types.StringValue(c.Type),
	}

	if c.Type == "PERSON" {
		obj, _ := types.ObjectValue(personAttrTypes, map[string]attr.Value{
			"first_name": types.StringValue(c.Person.FirstName),
			"last_name":  types.StringValue(c.Person.LastName),
			"email":      types.StringValue(c.Person.Email),
			"username":   types.StringValue(c.Person.Username),
			"department": types.StringValue(c.Person.Department),
			"group":      types.StringValue(c.Person.Group),
		})
		elems["person"] = obj
	} else {
		elems["person"] = types.ObjectNull(personAttrTypes)
	}

	if c.Type == "EGROUP" {
		obj, _ := types.ObjectValue(egroupAttrTypes, map[string]attr.Value{
			"name":  types.StringValue(c.EGroup.Name),
			"email": types.StringValue(c.EGroup.Email),
		})
		elems["egroup"] = obj
	} else {
		elems["egroup"] = types.ObjectNull(egroupAttrTypes)
	}

	if c.Type == "RESERVED" {
		obj, _ := types.ObjectValue(reservedAttrTypes, map[string]attr.Value{
			"first_name": types.StringValue(c.Reserved.FirstName),
			"last_name":  types.StringValue(c.Reserved.LastName),
		})
		elems["reserved"] = obj
	} else {
		elems["reserved"] = types.ObjectNull(reservedAttrTypes)
	}

	obj, _ := types.ObjectValue(contactAttrTypes, elems)
	return obj
}

func operatingSystemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"family":  types.StringType,
		"version": types.StringType,
	}
}

func flattenOperatingSystem(os landb.OperatingSystem) types.Object {
	return types.ObjectValueMust(
		operatingSystemAttrTypes(),
		map[string]attr.Value{
			"family":  types.StringValue(os.Family),
			"version": types.StringValue(os.Version),
		},
	)
}

func expandOperatingSystem(ctx context.Context, obj types.Object) (landb.OperatingSystem, diag.Diagnostics) {
	var diags diag.Diagnostics
	var out landb.OperatingSystem

	if obj.IsNull() || obj.IsUnknown() {
		return out, diags
	}

	attrs := obj.Attributes()

	var family types.String
	diags.Append(tfsdk.ValueAs(ctx, attrs["family"], &family)...)
	out.Family = family.ValueString()

	var version types.String
	diags.Append(tfsdk.ValueAs(ctx, attrs["version"], &version)...)
	out.Version = version.ValueString()

	out.Version = version.ValueString()

	return out, diags
}

func locationAttrTypes() map[string]attr.Type {
    return map[string]attr.Type{
        "building": types.StringType,
        "floor":    types.StringType,
        "room":     types.StringType,
    }
}

func flattenLocation(location landb.Location) types.Object {
    return types.ObjectValueMust(
        locationAttrTypes(),
        map[string]attr.Value{
            "building": types.StringValue(location.Building),
            "floor":    types.StringValue(location.Floor),
            "room":     types.StringValue(location.Room),
        },
    )
}

func expandLocation(ctx context.Context, obj types.Object) (landb.Location, diag.Diagnostics) {
    var diags diag.Diagnostics
    var out landb.Location

    if obj.IsNull() || obj.IsUnknown() {
        return out, diags
    }

    attrs := obj.Attributes()

    var b types.String
    diags.Append(tfsdk.ValueAs(ctx, attrs["building"], &b)...)
    out.Building = b.ValueString()

    var f types.String
    diags.Append(tfsdk.ValueAs(ctx, attrs["floor"], &f)...)
    out.Floor = f.ValueString()

    var r types.String
    diags.Append(tfsdk.ValueAs(ctx, attrs["room"], &r)...)
    out.Room = r.ValueString()

    return out, diags
}