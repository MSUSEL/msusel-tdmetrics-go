package typeDesc

type TypeDesc interface {
	Equal(other TypeDesc) bool
}

func listEqual[U TypeDesc](list1, list2 []U) bool {
	if len(list1) != len(list2) {
		return false
	}
	for i, elem := range list1 {
		if !elem.Equal(list2[i]) {
			return false
		}
	}
	return true
}
