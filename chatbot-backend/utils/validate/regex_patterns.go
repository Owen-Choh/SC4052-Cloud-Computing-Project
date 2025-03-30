package validate

import "regexp"

var ValidFileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\\-\\. ]+$`)