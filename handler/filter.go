package handler

import (
	"errors"
	"os"
	"strings"
)

func filter(email string, attr string) error {
	// Email domain filter
	if !validDomain(email) {
		return errors.New("Your email domain is not listed as authorized domain.")
	}
	return nil
}

func validDomain(e string) bool {
	var domains []string
	// Is the filter set?
	if f := os.Getenv("DOMAINFILTER"); f == "" {
		// If not set, consider everything valid
		return true
	} else {
		domains = strings.Split(f, ",")
	}
	// Get the domain out of the email
	s := strings.Split(e, "@")
	if len(s) != 2 {
		// Seems username is not an email address
		return false
	}
	for _, d := range domains {
		// If email domain matches with one of the authorized domains, its all good!
		if s[1] == d {
			return true
		}
	}
	return false
}
