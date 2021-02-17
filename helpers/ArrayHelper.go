package helpers

//TODO - why is this how we have to do it, horrible....
func StringArrayContainsItem(array []string, item string) bool {
	for _, arrayItem := range array {
		if item == arrayItem {
			return true
		}
	}
	return false
}
