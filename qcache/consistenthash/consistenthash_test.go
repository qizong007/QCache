package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 自定义hash函数
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 2, 12, 22
	// 4, 14, 24
	// 6, 16, 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2", // 2
		"11": "2", // 12
		"23": "4", // 24
		"27": "2", // 2
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 8, 18, 28
	hash.Add("8")

	testCases["27"] = "8" // 28

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
