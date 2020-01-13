package util

func MatchStringInSlice(tomatch string, sl []string) bool {
	for _, s := range sl {
		if s == tomatch {
			return true
		}
	}
	return false
}
