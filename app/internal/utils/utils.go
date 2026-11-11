package utils

func MapToString(input map[string]string) []string {
	outputStr := []string{}
	for key, value := range input {
		outputStr = append(outputStr, key, value)
	}
	return outputStr
}
