package chekers

import "encoding/json"

func ContainsInString(jsonString string, item string) (bool, error) {
	var slice []string
	err := json.Unmarshal([]byte(jsonString), &slice)
	if err != nil {
		return false, err
	}

	for _, s := range slice {
		if s == item {
			return true, nil
		}
	}
	return false, nil
}
