package algorithms

import (
	"fmt"
	"sort"
	"strings"
)

func ArePermutation(s1 string, s2 string) bool {
	// fmt.Println("arePermutation: str1=", s1, "str2=", s2)

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

func Generate_key(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	fmt.Println("Generate_key", strings.Join(s, ""))
	return strings.Join(s, "")
}
