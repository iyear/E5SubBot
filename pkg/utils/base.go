package utils

func IF(cond bool, a, b interface{}) interface{} {
	if cond {
		return a
	}
	return b
}
