package tools

import "strings"

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
