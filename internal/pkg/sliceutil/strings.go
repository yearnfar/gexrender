package sliceutil

// InStrings 是否在数组中
func InStrings(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
