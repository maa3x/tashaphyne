package tashaphyne

import (
	"regexp"
)

var Verbs map[string]struct{}

func init() {
	Verbs = make(map[string]struct{}, len(verbs))
	for i := range verbs {
		Verbs[VerbStamp(Runes(verbs[i]))] = struct{}{}
	}
}

func IsVerbStamp(s string) bool {
	_, ok := Verbs[s]
	return ok
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
