package validate

import "regexp"

var ValidChatbotNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\\-]+$`)
var ValidFileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\\-\\. ]+$`)