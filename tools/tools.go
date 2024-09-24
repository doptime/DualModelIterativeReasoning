package tools

import (
	"fmt"
	"strconv"
	"strings"
)

func TagContent(s string, tagSeperator string) (tag string) {
	ind := strings.Index(s, tagSeperator)
	if ind < 0 {
		return ""
	}
	endingTagSeperator := "</" + tagSeperator[1:]
	ind2 := strings.Index(s, endingTagSeperator)
	if ind2 < 0 {
		return ""
	}
	return s[ind+len(tagSeperator) : ind2]
}

func ReadFloatAfterTag(s string, tags ...string) (float64, error) {
	ind, tag := -1, ""
	for i := 0; i < len(tags) && ind < 0; i++ {
		tag = tags[i]
		ind = strings.Index(s, tag)
	}
	if ind < 0 {
		return 0, nil
	}
	s = s[ind+len(tag):]
	s = strings.TrimSpace(s)

	s = strings.Split(s, "\n")[0]
	if ind := strings.Index(s, "="); ind >= 0 {
		s = s[ind+1:]
	}
	s = strings.TrimSpace(s)

	validInd := 0
	for ; validInd < len(s) && strings.ContainsRune("0123456789.", rune(s[validInd])); validInd++ {
	}
	s = s[:validInd]
	if validInd == 0 {
		return 0, fmt.Errorf("no number found")
	}
	return strconv.ParseFloat(s, 64)
}
