package chord

const (
	PatternSize int = 100
)

// Search Last Position String
func SearchNext(haystack string, needle string) int {
	retSlice := KMP(haystack, needle)
	if len(retSlice) > 0 {
		return retSlice[len(retSlice)-1]
	}

	return -1
}

// Search First Position String
func SearchString(haystack string, needle string) int {
	retSlice := KMP(haystack, needle)
	if len(retSlice) > 0 {
		return retSlice[0]
	}

	return -1
}

func KMP(haystack string, needle string) []int {
	next := preKMP(needle)
	i := 0
	j := 0
	m := len(needle)
	n := len(haystack)

	x := []byte(needle)
	y := []byte(haystack)
	var ret []int

	// got zero target or want, just return empty result
	if m == 0 || n == 0 {
		return ret
	}

	// want string bigger than target string
	if n < m {
		return ret
	}

	for j < n {
		for i > -1 && x[i] != y[j] {
			i = next[i]
		}
		i++
		j++

		if i >= m {
			ret = append(ret, j-i)
			i = next[i]
		}
	}

	return ret
}

func preKMP(x string) [PatternSize]int {
	var i, j int
	length := len(x) - 1
	var kmpNext [PatternSize]int
	i = 0
	j = -1
	kmpNext[0] = -1

	for i < length {
		for j > -1 && x[i] != x[j] {
			j = kmpNext[j]
		}

		i++
		j++

		if x[i] == x[j] {
			kmpNext[i] = kmpNext[j]
		} else {
			kmpNext[i] = j
		}
	}
	return kmpNext
}