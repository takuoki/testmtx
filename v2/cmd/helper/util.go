package helper

import "fmt"

func appendUsageItem(usage, item string, i, length int) string {
	if length < 2 {
		panic("length must be greater than or equal to 2")
	}

	if i == 0 {
		return fmt.Sprintf("%s (%q", usage, item)
	}
	if i < length-1 {
		return fmt.Sprintf("%s, %q", usage, item)
	}
	return fmt.Sprintf("%s or %q)", usage, item)
}
