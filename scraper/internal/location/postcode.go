package location

import (
	"errors"
	"regexp"
	"strings"
)

var errNoPostcode = errors.New("could not get postcode")

func getPostcodeFromAddress(address string) (string, error) {
	parts := strings.Split(address, ",")
	for _, part := range parts {
		val := strings.TrimSpace(part)
		_isPostcode, err := isPostcode(val)
		if err != nil {
			return "", err
		}

		if _isPostcode {
			return val, nil
		}
	}

	return "", errNoPostcode
}

// https://stackoverflow.com/questions/164979/regex-for-matching-uk-postcodes#164994
const postcodePattern = "^([A-Za-z][A-Ha-hJ-Yj-y]?[0-9][A-Za-z0-9]? ?[0-9][A-Za-z]{2}|[Gg][Ii][Rr] ?0[Aa]{2})$"

func isPostcode(val string) (bool, error) {
	return regexp.Match(postcodePattern, []byte(val))
}
