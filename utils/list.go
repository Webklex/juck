package utils

func UniqueStringList(list []string) (uniqueList []string) {
	if list == nil {
		return
	}
	for _, s := range list {
		if InStringList(uniqueList, s) == false {
			uniqueList = append(uniqueList, s)
		}
	}
	return
}

func InStringList(list []string, target string) bool {
	if list == nil {
		return false
	}
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
