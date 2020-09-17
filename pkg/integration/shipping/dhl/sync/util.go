package sync

import "o.o/capi/dot"

func max(a, b dot.ID) dot.ID {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func convertIDsToStrings(ids []dot.ID) []string {
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, id.String())
	}
	return idsStr
}
