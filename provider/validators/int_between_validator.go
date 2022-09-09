package validators

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func IntBetween(min, max int64) tfsdk.AttributeValidator {
	return IntBetweenValidator{Min: min, Max: max}
}

func IntGreaterThan(min int64) tfsdk.AttributeValidator {
	return IntBetweenValidator{Min: min, Max: math.MaxInt64}
}

func IntSmallerThan(max int64) tfsdk.AttributeValidator {
	return IntBetweenValidator{Min: math.MinInt64, Max: max}
}

type IntBetweenValidator struct {
	Max int64
	Min int64
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v IntBetweenValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("value must be between `%d` and `%d`", v.Min, v.Max)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v IntBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("value must be between `%d` and `%d`", v.Min, v.Max)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v IntBetweenValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic validators, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var value types.Int64
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &value)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if value.Unknown || value.Null {
		return
	}

	if value.Value < v.Min || value.Value > v.Max {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Numeric Value",
			fmt.Sprintf("value must be between %d and %d", v.Min, v.Max),
		)
	}
}
