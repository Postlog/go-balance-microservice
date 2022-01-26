package utils

func StringInCollection(target string, col ...string) bool {
	for _, s := range col {
		if s == target {
			return true
		}
	}
	return false
}
