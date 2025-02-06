package tashaphyne

import (
	"regexp"
	"slices"
	"strings"
)

type Runes []rune

func (r Runes) String() string {
	return string(r)
}

func (r Runes) Replace(old, new string) Runes {
	return Runes(strings.ReplaceAll(string(r), old, new))
}

func (r Runes) ReplaceRegex(re *regexp.Regexp, new string) Runes {
	return Runes(re.ReplaceAllString(string(r), new))
}

func (r Runes) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(r), prefix)
}

func (r Runes) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(r), suffix)
}

func (r Runes) Contains(sub string) bool {
	return strings.Contains(string(r), sub)
}

func (r Runes) In(values ...string) bool {
	return slices.Contains(values, string(r))
}

func (r Runes) StripTashkeel() Runes {
	return Runes(StripTashkeel(string(r)))
}

func (r Runes) At(pos int) string {
	if pos < 0 {
		pos = len(r) + pos
	}
	if pos < 0 || pos >= len(r) {
		return ""
	}
	return string(r[pos])
}

func (r Runes) Slice(from, to int) Runes {
	if from < 0 {
		from = len(r) + from
	}
	if to < 0 {
		to = len(r) + to
	}

	if from < 0 || from > len(r) || to < 0 || to > len(r) || from > to {
		return Runes{}
	}

	return r[from:to]
}

func (r Runes) SliceFrom(from int) Runes {
	if from < 0 {
		from = len(r) + from
	}
	if from < 0 || from > len(r) {
		return Runes{}
	}
	return r[from:]
}

func (r Runes) SliceTo(to int) Runes {
	if to < 0 {
		to = len(r) + to
	}
	if to < 0 || to > len(r) {
		return Runes{}
	}
	return r[:to]
}

func (r Runes) IndexOf(s string) int {
	for i, v := range r {
		if string(v) == s {
			return i
		}
	}
	return -1
}

func (r Runes) LastIndexOf(s string) int {
	for i := len(r) - 1; i >= 0; i-- {
		if string(r[i]) == s {
			return i
		}
	}
	return -1
}

func (r Runes) Append(s string) Runes {
	return append(r, []rune(s)...)
}

func (r Runes) Prepend(s string) Runes {
	return append([]rune(s), r...)
}
