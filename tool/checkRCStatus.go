package tool

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zenmuharom/zenlogger"
)

func CheckRCStatus(logger zenlogger.Zenlogger, rc string, rcs []uint8) (valid bool) {
	rcsString := fmt.Sprintf("%s", rcs)
	logger.Debug("CheckRCStatus", zenlogger.ZenField{Key: "rc", Value: rc}, zenlogger.ZenField{Key: "rcs", Value: rcsString})

	// Create a regular expression pattern to match string values
	pattern := `"[^"]+"`

	// Compile the regex pattern
	re := regexp.MustCompile(pattern)

	// Find all matches of integers in the rcs string
	matches := re.FindAllString(rcsString, -1)

	// Remove the double quotes from the matches
	for i, match := range matches {
		matches[i] = strings.Trim(match, `"`)
	}

	// Check if rc value exists in matches
	for _, match := range matches {
		if match == rc {
			logger.Debug("checkRCStatus", zenlogger.ZenField{Key: "compare", Value: fmt.Sprintf("(match) %v == (rc) %v", match, rc)})
			valid = true
			break
		}
	}
	return
}
