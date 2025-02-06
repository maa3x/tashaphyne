package tashaphyne

import (
	"regexp"
	"slices"
)

var VerbList []string

func init() {
	VerbList = make([]string, len(verbs))
	for i := range verbs {
		VerbList[i] = VerbStamp(Runes(verbs[i]))
	}
}

func IsVerbStamp(s string) bool {
	return slices.Contains(VerbList, s)
}

func VerbStamp(s Runes) string {
	if len(s) == 0 {
		return ""
	}
	s = s.StripTashkeel()
	if len(s) == 0 {
		return ""
	}

	if len(s) == 4 && s.At(0) == AlefHamzaAbove {
		s = s.Replace(AlefHamzaAbove, "")
	}
	s = Runes(regexp.MustCompile("[أإءؤئآ]").ReplaceAllString(string(s), "ء"))
	s = Runes(regexp.MustCompile("[اويى]").ReplaceAllString(string(s), ""))
	if string(s.SliceFrom(-1)) == string(s.Slice(-2, -1)) {
		s = s.SliceTo(-1)
	}

	return string(s)
}
