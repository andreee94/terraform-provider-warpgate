package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringIsDomain() tfsdk.AttributeValidator {
	return StringRegex{
		Regex: regexp.MustCompile(`^((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))$`),
	}
}

func StringIsIpv4() tfsdk.AttributeValidator {
	return StringRegex{
		Regex: regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`),
	}
}

func StringIsMacAddress() tfsdk.AttributeValidator {
	return StringRegex{
		Regex: regexp.MustCompile(`^[a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}$`),
	}
}

type StringRegex struct {
	Regex *regexp.Regexp
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringRegex) Description(ctx context.Context) string {
	return fmt.Sprintf("string must match regex %s", v.Regex.String())
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringRegex) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string must match regex `%s`", v.Regex.String())
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v StringRegex) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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

	if !v.Regex.MatchString(str.Value) {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String",
			fmt.Sprintf("string must match regex %s", v.Regex.String()),
		)

		return
	}
}
