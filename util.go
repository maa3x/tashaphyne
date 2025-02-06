package tashaphyne

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
)

func Min[T cmp.Ordered](values ...T) (m T) {
	for i := range values {
		if i == 0 {
			m = values[i]
		} else if values[i] < m {
			m = values[i]
		}
	}

	return m
}

func Max[T cmp.Ordered](values ...T) (m T) {
	for i := range values {
		if i == 0 {
			m = values[i]
		} else if values[i] > m {
			m = values[i]
		}
	}

	return m
}

func In[T comparable](needle T, haystack ...T) bool {
	return slices.Contains(haystack, needle)
}

func re(format string, args ...any) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(format, args...))
}
