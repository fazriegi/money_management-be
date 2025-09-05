package libs

func Intersection[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]struct{})
	for _, v := range slice1 {
		set[v] = struct{}{}
	}

	var matches []T
	for _, v := range slice2 {
		if _, ok := set[v]; ok {
			matches = append(matches, v)
		}
	}

	return matches
}
