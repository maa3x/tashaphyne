package tashaphyne

const (
	Comma          = "\u060C"
	Semicolon      = "\u061B"
	Question       = "\u061F"
	Hamza          = "\u0621"
	AlefMadda      = "\u0622"
	AlefHamzaAbove = "\u0623"
	WawHamza       = "\u0624"
	AlefHamzaBelow = "\u0625"
	YehHamza       = "\u0626"
	Alef           = "\u0627"
	Beh            = "\u0628"
	TehMarbuta     = "\u0629"
	Teh            = "\u062a"
	Theh           = "\u062b"
	Jeem           = "\u062c"
	Hah            = "\u062d"
	Khah           = "\u062e"
	Dal            = "\u062f"
	Thal           = "\u0630"
	Reh            = "\u0631"
	Zain           = "\u0632"
	Seen           = "\u0633"
	Sheen          = "\u0634"
	Sad            = "\u0635"
	Dad            = "\u0636"
	Tah            = "\u0637"
	Zah            = "\u0638"
	Ain            = "\u0639"
	Ghain          = "\u063a"
	Tatweel        = "\u0640"
	Feh            = "\u0641"
	Qaf            = "\u0642"
	Kaf            = "\u0643"
	Lam            = "\u0644"
	Meem           = "\u0645"
	Noon           = "\u0646"
	Heh            = "\u0647"
	Waw            = "\u0648"
	AlefMaksura    = "\u0649"
	Yeh            = "\u064a"
	MaddaAbove     = "\u0653"
	HamzaAbove     = "\u0654"
	HamzaBelow     = "\u0655"
	Zero           = "\u0660"
	One            = "\u0661"
	Two            = "\u0662"
	Three          = "\u0663"
	Four           = "\u0664"
	Five           = "\u0665"
	Six            = "\u0666"
	Seven          = "\u0667"
	Eight          = "\u0668"
	Nine           = "\u0669"
	Percent        = "\u066a"
	Decimal        = "\u066b"
	Thousands      = "\u066c"
	Star           = "\u066d"
	MiniAlef       = "\u0670"
	AlefWasla      = "\u0671"
	FullStop       = "\u06d4"
	ByteOrderMark  = "\ufeff"

	// Diacritics
	Fathatan = "\u064b"
	Dammatan = "\u064c"
	Kasratan = "\u064d"
	Fatha    = "\u064e"
	Damma    = "\u064f"
	Kasra    = "\u0650"
	Shadda   = "\u0651"
	Sukun    = "\u0652"

	// Ligatures
	LamAlef                 = "\ufefb"
	LamAlefHamzaAbove       = "\ufef7"
	LamAlefHamzaBelow       = "\ufef9"
	LamAlefMaddaAbove       = "\ufef5"
	SimpleLamAlef           = "\u0644\u0627"
	SimpleLamAlefHamzaAbove = "\u0644\u0623"
	SimpleLamAlefHamzaBelow = "\u0644\u0625"
	SimpleLamAlefMaddaAbove = "\u0644\u0622"
)

var (
	harakatRegex   = re("[%s]", Fathatan+Dammatan+Kasratan+Fatha+Damma+Kasra+Sukun+Shadda)
	hamzatRegex    = re("[%s]", WawHamza+YehHamza)
	alefatRegex    = re("[%s]", AlefMadda+AlefHamzaAbove+AlefHamzaBelow+HamzaAbove+HamzaBelow)
	lamalefatRegex = re("[%s]", LamAlef+LamAlefHamzaAbove+LamAlefHamzaBelow+LamAlefMaddaAbove)
	tatweelRegex   = re("[%s]", Tatweel)
	alefMaddaRegex = re("[%s]", AlefMadda)
	tehTahDalRegex = re("[%s]", Teh+Tah+Dal)
	tehRegex       = re(Teh)
	tahRegex       = re(Tah)
	dalRegex       = re(Dal)
)
