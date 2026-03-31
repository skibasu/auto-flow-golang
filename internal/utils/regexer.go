package utils

import "regexp"

var (
	UpperRegex   = regexp.MustCompile(`[A-Z]`)
	DigitRegex   = regexp.MustCompile(`\d`)
	SpecialRegex = regexp.MustCompile(`[!@#$%^&*()_\-+=<>?{}\[\]~]`)
)
