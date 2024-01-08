package helper

import "strings"

func SplitComma(val string) []string {
	split := strings.Split(val, ",")

	if len(split) <= 0 {
		return []string{}
	} else {
		return split
	}
}
