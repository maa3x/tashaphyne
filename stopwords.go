package tashaphyne

func IsStop(s string) bool {
	_, ok := StopWords[s]
	return ok
}

func StopStem(s string) string {
	if stop, ok := StopWords[s]; ok {
		return StripTashkeel(stop.Stem)
	}
	return ""
}
