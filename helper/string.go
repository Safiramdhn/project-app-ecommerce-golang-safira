package helper

func JoinStrings(parts []string, delimiter string) string {
	result := ""
	for i, part := range parts {
		result += part
		if i < len(parts)-1 {
			result += delimiter
		}
	}
	return result
}
