package cherrySlice

func StringIn(v string, sl []string) (int, bool) {
	for i, s := range sl {
		if s == v {
			return i, true
		}
	}

	// TODO 返回 -1 最好
	return 0, false
}

func StringInSlice(v string, sl []string) bool {
	_, ok := StringIn(v, sl)
	return ok
}
