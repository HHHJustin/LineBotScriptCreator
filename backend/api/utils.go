package api

func getColorByType(nodeType string) string {
	switch nodeType {
	case "FirstStep":
		return "lightred"
	case "Message":
		return "lightgreen"
	case "QuickReply":
		return "yellow"
	case "KeywordDecision":
		return "orange"
	case "TagDecision":
		return "brown"
	case "TagOperation":
		return "pink"
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
