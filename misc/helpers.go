package misc

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func deleteFromSlice(list []string, a string) []string {
	for s, b := range list {
		if b == a {
			list = append(list[:s], list[s+1:]...)
			break
		}
	}
	return list
}
