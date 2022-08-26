package algorithms

import (
	"fmt"
)

func arePermutation(s1 string, s2 string) bool {
	fmt.Printf("arePermutation: str1=%s, str2=%s", s1, s2)

	map_counter := make(map[byte]int)

	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		map_counter[s1[i]]++
		map_counter[s2[i]]--

	}

	for i := 0; i < len(s1); i++ {
		if i != 0 {
			return false
		}
	}

	for _, value := range map_counter {
		if value != 0 {
			return false
		}
	}
	return true
}
