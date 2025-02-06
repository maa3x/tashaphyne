package tashaphyne

import "regexp"

func StripTashkeel(s string) string {
	return harakatRegex.ReplaceAllString(s, "")
}

func StripTatweel(s string) string {
	return tatweelRegex.ReplaceAllString(s, "")
}

func NormalizeHamza(s string) string {
	s = alefatRegex.ReplaceAllString(s, Alef)
	return hamzatRegex.ReplaceAllString(s, Hamza)
}

func NormalizeLamAlef(s string) string {
	return lamalefatRegex.ReplaceAllString(s, Lam+Alef)
}

func NormalizeSpellErrors(s string) string {
	s = regexp.MustCompile("["+TehMarbuta+"]").ReplaceAllString(s, Heh)
	return regexp.MustCompile("["+AlefMaksura+"]").ReplaceAllString(s, Yeh)
}

func NormalizeSearchText(s string) string {
	s = StripTashkeel(s)
	s = StripTatweel(s)
	s = NormalizeLamAlef(s)
	s = NormalizeHamza(s)
	s = NormalizeSpellErrors(s)
	return s
}
