package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringIn(validValues []string, ignoreCase bool) tfsdk.AttributeValidator {
	return StringInValidator{
		ValidValues: validValues,
		IgnoreCase:  ignoreCase,
	}
}

type StringInValidator struct {
	ValidValues []string
	IgnoreCase  bool
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringInValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("String must be one of %v.", v.ValidValues)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringInValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("String must be one of `%v`.", v.ValidValues)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v StringInValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic validators, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.Unknown || str.Null {
		return
	}

	// probably is better to reverse for loop and if case, but we do not really care too much about performance here
	for _, validValue := range v.ValidValues {
		if v.IgnoreCase {
			if strings.ToLower(str.Value) == strings.ToLower(validValue) {
				return
			}
		} else {
			if str.Value == validValue {
				return
			}
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Invalid String",
		fmt.Sprintf("String must be one of %v, got: %s.", v.ValidValues, str.Value),
	)
	return
}
