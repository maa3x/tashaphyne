package tashaphyne

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

type Stemmer struct {
	prefixLetters        string
	suffixLetters        string
	infixLetters         string
	maxPrefixLength      int
	maxSuffixLength      int
	minStemLength        int
	joker                string
	prefixList           []string
	suffixList           []string
	rootList             []string
	validAffixesList     []string
	word                 Runes
	starWord             Runes
	unvocalized          Runes
	normalized           Runes
	root                 Runes
	left                 int
	right                int
	stemLeft             int
	stemRight            int
	segmentList          [][2]int
	prefixesTree         map[string]any
	suffixesTree         map[string]any
	tokenRegex           *regexp.Regexp
	notPrefixRegex       *regexp.Regexp
	notInfixRegex        *regexp.Regexp
	notSuffixRegex       *regexp.Regexp
	notPrefixSuffixRegex *regexp.Regexp
}

func New() *Stemmer {
	s := &Stemmer{
		prefixLetters:    DefaultPrefixLetters,
		suffixLetters:    DefaultSuffixLetters,
		infixLetters:     DefaultInfixLetters,
		maxPrefixLength:  DefaultMaxPrefix,
		maxSuffixLength:  DefaultMaxSuffix,
		minStemLength:    DefaultMinStem,
		joker:            DefaultJoker,
		prefixList:       DefaultPrefixList,
		suffixList:       DefaultSuffixList,
		rootList:         Roots,
		validAffixesList: AffixList,
		tokenRegex:       re("[^\\w\u064b-\u0652']+"),
	}
	s.notPrefixRegex = re("[^%s]", s.prefixLetters)
	s.notInfixRegex = re("[^%s]", s.infixLetters)
	s.notSuffixRegex = re("[^%s]", s.suffixLetters)
	s.notPrefixSuffixRegex = re("[^%s]", s.prefixLetters+s.suffixLetters)
	s.prefixesTree = s.createPrefixesTree(s.prefixList)
	s.suffixesTree = s.createSuffixesTree(s.suffixList)

	return s
}

func (s *Stemmer) LightStem(word string) string {
	if word == "" {
		return ""
	}

	s.transformToStars(Runes(word))
	s.segment(Runes(word))
	return string(s.getStem(-1, -1))
}

func (s *Stemmer) createPrefixesTree(prefixes []string) map[string]any {
	prefixesTree := make(map[string]any)

	for _, prefix := range prefixes {
		branch := prefixesTree
		for _, ch := range Runes(prefix) {
			char := string(ch)
			if _, exists := branch[char]; !exists {
				branch[char] = make(map[string]any)
			}
			branch = branch[char].(map[string]any)
		}

		if endMarker, exists := branch["#"]; exists {
			endMarker.(map[string]any)[prefix] = "#"
		} else {
			branch["#"] = map[string]any{prefix: "#"}
		}
	}

	return prefixesTree
}

func (s *Stemmer) createSuffixesTree(suffixes []string) map[string]any {
	suffixesTree := make(map[string]any)
	for _, suffix := range suffixes {
		branch := suffixesTree
		for _, ch := range Runes(suffix) {
			char := string(ch)
			if _, exists := branch[char]; !exists {
				branch[char] = make(map[string]any)
			}
			branch = branch[char].(map[string]any)
		}

		if endMarker, exists := branch["#"]; exists {
			endMarker.(map[string]any)[suffix] = "#"
		} else {
			branch["#"] = map[string]any{suffix: "#"}
		}
	}

	return suffixesTree
}

func (s *Stemmer) getRoot(preIdx, suIdx int) Runes {
	if preIdx >= 0 || suIdx >= 0 {
		s.extractRoot(preIdx, suIdx)
	} else {
		s.root = s.chooseRoot()
	}

	return s.root
}

func (s *Stemmer) chooseRoot() Runes {
	if IsStop(string(s.word)) {
		return Runes(StopStem(string(s.word)))
	}

	if len(s.segmentList) == 0 {
		s.segment(s.word)
	}
	affixList := s.getAffixList(nil)
	var roots []string
	for i := range affixList {
		roots = append(roots, string(affixList[i].root))
	}

	rootsTmp := roots
	var accepted []string
	for i := range rootsTmp {
		if s.isRootLengthValid(Runes(rootsTmp[i])) {
			accepted = append(accepted, rootsTmp[i])
		}
	}
	if len(accepted) > 0 {
		rootsTmp = accepted
	}

	accepted = nil
	for i := range rootsTmp {
		if s.isRoot(Runes(rootsTmp[i])) {
			accepted = append(accepted, rootsTmp[i])
		}
	}
	if len(accepted) > 0 {
		rootsTmp = accepted
	}

	return s.mostCommon(rootsTmp)
}

func (s *Stemmer) chooseStem() Runes {
	if IsStop(string(s.word)) {
		return Runes(StopStem(string(s.word)))
	}

	if len(s.segmentList) == 0 {
		s.segment(s.word)
	}
	var segList [][2]int
	for i := range s.segmentList {
		if s.verifyAffix(s.segmentList[i][0], s.segmentList[i][1]) {
			segList = append(segList, s.segmentList[i])
		}
	}

	if len(segList) == 0 {
		return s.unvocalized
	}

	left, right := s.getLeftRight(segList)
	return s.unvocalized.Slice(left, right)
}

func (s *Stemmer) getStem(preIdx, suIdx int) Runes {
	var left, right int
	if preIdx >= 0 || suIdx >= 0 {
		if preIdx > 0 {
			left = preIdx
		} else {
			left = s.stemLeft
		}
		if suIdx > 0 {
			right = suIdx
		} else {
			right = s.stemRight
		}
		return s.unvocalized.Slice(left, right)
	}

	return s.chooseStem()
}

func (s *Stemmer) handleTehInflix(word Runes, left, right int) Runes {
	if keyStem := word.Replace(TehMarbuta, ""); len(keyStem) != 4 {
		return word.ReplaceRegex(tehTahDalRegex, s.joker)
	}

	word = append(word[:2], word[2:].ReplaceRegex(tehRegex, s.joker)...)
	if s.word.Slice(left, right).HasPrefix("ضط") {
		word = append(word[:2], word[2:].ReplaceRegex(tahRegex, s.joker)...)
	} else {
		word = word.ReplaceRegex(tahRegex, s.joker)
	}

	if s.word[left:right].HasPrefix("زد") {
		word = append(word[:2], word[2:].ReplaceRegex(dalRegex, s.joker)...)
	} else {
		word = word.ReplaceRegex(dalRegex, s.joker)
	}

	return word
}

func (s *Stemmer) getStarStem(preIdx, suIdx int) Runes {
	word := s.word
	if preIdx < 0 && suIdx < 0 {
		return word.Slice(s.left, s.right)
	}

	left, right := s.left, s.right
	if preIdx >= 0 {
		left = preIdx
	}
	if suIdx >= 0 {
		right = suIdx
	}

	if s.infixLetters != "" {
		return s.handleTehInflix(word.Slice(left, right).ReplaceRegex(re("[^%s]", s.infixLetters+TehMarbuta), s.joker), left, right)
	}

	return Runes(strings.Repeat(s.joker, len(word.Slice(left, right))))
}

func (s *Stemmer) getPrefix(preIdx int) Runes {
	if preIdx >= 0 {
		return s.unvocalized.SliceTo(preIdx)
	}
	return s.unvocalized.SliceTo(s.left)
}

func (s *Stemmer) getSuffix(suIdx int) Runes {
	if suIdx >= 0 {
		return s.unvocalized.SliceFrom(suIdx)
	}
	return s.unvocalized.SliceFrom(s.right)
}

func (s *Stemmer) getAffix(preIdx, suIdx int) Runes {
	return Runes(string(s.getPrefix(preIdx)) + "-" + string(s.getSuffix(suIdx)))
}

type affixTuple struct {
	prefix   Runes
	suffix   Runes
	stem     Runes
	starStem Runes
	root     Runes
}

func (s *Stemmer) getAffixTuple(preIdx, suIdx int) affixTuple {
	return affixTuple{
		prefix:   s.getPrefix(preIdx),
		suffix:   s.getSuffix(suIdx),
		stem:     s.getStem(preIdx, suIdx),
		starStem: s.getStarStem(preIdx, suIdx),
		root:     s.getRoot(preIdx, suIdx),
	}
}

func (s *Stemmer) transformToStars(word Runes) (Runes, int, int) {
	s.word = word
	word = word.StripTashkeel()
	s.unvocalized = word
	word = word.ReplaceRegex(alefMaddaRegex, Hamza+Alef)
	word = word.ReplaceRegex(s.notPrefixSuffixRegex, s.joker)

	left := word.IndexOf(s.joker)
	right := word.LastIndexOf(s.joker)
	if left >= 0 {
		left = Min(left, s.maxPrefixLength-1)
		right = max(right+1, len(word)-s.maxSuffixLength)
		prefix := word.SliceTo(left)
		stem := s.word.Slice(left, right)
		suffix := word.SliceFrom(right)
		if s.infixLetters != "" {
			stem = stem.ReplaceRegex(s.notInfixRegex, s.joker)
		}
		word = Runes(string(prefix) + string(stem) + string(suffix))
	}

	left = word.IndexOf(s.joker)
	right = word.LastIndexOf(s.joker)
	if left < 0 {
		left = Min(s.maxPrefixLength, len(word)-2)
	}
	if left >= 0 {
		prefix := word.SliceTo(left)
		for len(prefix) > 0 && !prefix.In(s.prefixList...) {
			prefix = prefix.SliceTo(-1)
		}
		if right < 0 {
			right = max(len(prefix), len(word)-s.maxSuffixLength)
		}
		suffix := word.SliceFrom(right)
		for len(suffix) > 0 && !suffix.In(s.suffixList...) {
			suffix = suffix.SliceFrom(1)
		}

		left = len(prefix)
		right = len(word) - len(suffix)
		stem := s.word.Slice(left, right)
		if s.infixLetters != "" {
			stem = stem.ReplaceRegex(s.notInfixRegex, s.joker)
		}
		word = Runes(string(prefix) + string(stem) + string(suffix))
	}

	s.stemLeft = left
	s.stemRight = right
	s.starWord = word

	return word, left, right
}

func (s *Stemmer) extractRoot(preIdx, suIdx int) Runes {
	stem := s.getStem(preIdx, suIdx)
	root := Runes("")
	if len(stem) == 3 {
		s.root = s.adjustRoot(root, stem)
		return s.root
	}

	starStem := s.getStarStem(preIdx, suIdx)
	if len(starStem) == len(stem) {
		for i, ch := range stem {
			if string(starStem[i]) == s.joker {
				root = append(root, ch)
			}
		}
	} else {
		root = stem
	}

	root = s.normalizeRoot(root)
	if len(root) == 2 {
		root = s.adjustRoot(root, starStem)
	}
	s.root = root
	return root
}

func (s *Stemmer) adjustRoot(root, starStem Runes) Runes {
	if len(starStem) == 0 {
		return root
	}

	if len(starStem) == 3 {
		return starStem.Replace(Alef, Waw).Replace(AlefMaksura, Yeh)
	}

	first := starStem.At(0)
	last := starStem.At(-1)
	if first == Alef || first == Waw {
		root = root.Prepend(Waw)
	} else if first == Yeh {
		root = root.Prepend(Yeh)
	} else if first == s.joker && In(last, Alef, Waw) {
		root = root.Append(Waw)
	} else if first == s.joker && In(last, AlefMaksura, Yeh) {
		root = root.Append(Waw)
	} else if first == s.joker && last == s.joker {
		if len(starStem) == 2 {
			root = root.Append(root.At(-1))
		} else {
			root = Runes(root.At(0) + Waw + root.At(1))
		}
	}

	return root
}

func (s *Stemmer) lookupPrefixes(word Runes) []int {
	branch := s.prefixesTree
	lefts := []int{0}

	i := 0
	for ; i < len(word); i++ {
		ch := word.At(i)
		if _, ok := branch[ch]; !ok {
			break
		}

		if _, hashOk := branch["#"]; hashOk {
			lefts = append(lefts, i)
		}
		branch = branch[ch].(map[string]any)
	}

	if _, hashExist := branch["#"]; hashExist && i < len(word) {
		lefts = append(lefts, i)
	}

	return lefts
}

func (s *Stemmer) lookupSuffixes(word Runes) []int {
	branch := s.suffixesTree
	suffix := Runes("")
	rights := []int{}

	i := len(word) - 1
	for ; i >= 0; i-- {
		ch := word.At(i)
		if _, ok := branch[ch]; !ok {
			break
		}

		suffix = suffix.Prepend(ch)
		if _, hashOk := branch["#"]; hashOk {
			rights = append(rights, i+1)
		}
		branch = branch[ch].(map[string]any)
	}

	if _, hashExist := branch["#"]; hashExist && i >= 0 {
		rights = append(rights, i+1)
	}

	return rights
}

func (s *Stemmer) segment(word Runes) [][2]int {
	s.word = word
	s.unvocalized = word.StripTashkeel()
	word = word.ReplaceRegex(alefMaddaRegex, Hamza+Alef)
	lefts := s.lookupPrefixes(word)
	rights := s.lookupSuffixes(word)
	if len(lefts) > 0 {
		s.left = Max(lefts...)
	} else {
		s.left = -1
	}

	if len(rights) > 0 {
		s.right = Min(rights...)
	} else {
		s.right = -1
	}

	s.segmentList = [][2]int{{0, len(word)}}
	for _, i := range lefts {
		for _, j := range rights {
			if j >= i+2 {
				s.segmentList = append(s.segmentList, [2]int{i, j})
			}
		}
	}

	s.left, s.right = s.getLeftRight(s.segmentList)
	return s.segmentList
}

func (s *Stemmer) getAffixList(segList [][2]int) []affixTuple {
	if len(segList) == 0 {
		segList = s.segmentList
	}
	var affixList []affixTuple
	for i := range segList {
		affixList = append(affixList, s.getAffixTuple(segList[i][0], segList[i][1]))
	}
	return affixList
}

func (s *Stemmer) validStem(stem Runes, tag string, prefix Runes) bool {
	if len(stem) == 0 {
		return false
	}

	if tag == Noun {
		return len(stem) < 8
	}

	// verb
	if len(stem) > 6 || len(stem) < 2 {
		return false
	} else if stem.Contains(TehMarbuta) {
		return false
	} else if len(stem) == 6 && !stem.HasPrefix(Alef) {
		return false
	} else if len(stem) == 5 && !(stem.At(0) == Alef || stem.At(0) == Teh) {
		if prefix.SliceFrom(-1).In(Yeh, Teh, Noon, AlefHamzaAbove) {
			return false
		}
	} else if stem.HasPrefix(Alef) && prefix.SliceFrom(-1).In(Yeh, Noon, Teh, AlefHamzaAbove, Alef) {
		return false
	}

	if !IsVerbStamp(string(stem)) {
		return false
	}
	return true
}

func (s *Stemmer) verifyAffix(preIdx, suIdx int) bool {
	prefix := s.getPrefix(preIdx)
	suffix := s.getSuffix(suIdx)
	affix := string(prefix) + "-" + string(suffix)
	stem := s.getStem(preIdx, suIdx)
	if slices.Contains(VerbAffixList, affix) && s.validStem(stem, Verb, prefix) {
		if slices.Contains(NounAffixList, affix) && s.validStem(stem, Noun, Runes("")) {
			return true
		}
		return true
	}
	if slices.Contains(NounAffixList, affix) && s.validStem(stem, Noun, Runes("")) {
		return true
	}
	return false
}

func (s *Stemmer) normalizeRoot(word Runes) Runes {
	word = word.Replace(AlefMadda, Hamza+Alef)
	word = word.Replace(TehMarbuta, "")
	word = word.Replace(AlefMaksura, Yeh)
	return Runes(NormalizeHamza(string(word)))
}

func (s *Stemmer) isRootLengthValid(root Runes) bool {
	return len(root) >= 2 && len(root) <= 4
}

func (s *Stemmer) mostCommon(lst []string) Runes {
	var triroots []string
	for i := range lst {
		if len(lst[i]) == 3 {
			triroots = append(triroots, lst[i])
		}
	}
	if len(triroots) > 0 {
		lst = triroots
	}

	return Runes(slices.MaxFunc(lst, func(a, b string) int { return cmp.Compare(len(a), len(b)) }))
}

func (s *Stemmer) isRoot(word Runes) bool {
	return slices.Contains(s.rootList, string(word))
}

func (s *Stemmer) getLeftRight(segList [][2]int) (int, int) {
	if len(segList) == 0 {
		return -1, -1
	}

	var l, r int
	for i := range segList {
		if i == 0 {
			l = segList[i][0]
		} else if segList[i][0] > l {
			l = segList[i][0]
		}
	}
	for i := range segList {
		if i == 0 {
			r = segList[i][1]
		} else if segList[i][0] == l && segList[i][1] < r {
			r = segList[i][1]
		}
	}

	return l, r
}
