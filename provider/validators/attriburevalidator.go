package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func IsDomain() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))$`),
		"Invalid domain",
	)
}

func IsIpv4() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`),
		"Invalid ipv4",
	)
}

func IsMacAddress() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^[a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}$`),
		"Invalid mac address",
	)
}

func IsUUID() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`),
		"Invalid uuid",
	)
}
