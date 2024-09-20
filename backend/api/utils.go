package api

func getColorByType(nodeType string) string {
	switch nodeType {
	case "start":
		return "lightgreen"
	case "message":
		return "lightgreen"
	case "quickreply":
		return "yellow"
	case "process":
		return "lightblue"
	default:
		return "lightgray"
	}
}

func contains(arr []int, value int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func removeValue(arr []int, value int) []int {
	var result []int
	for _, v := range arr {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
