package currencyn

const (
	TRC = "TRC" // 体验币

	multi     = 1000
	precision = 3
)

type Conf struct {
	Currency  string `json:"currency"`  // 货币代码
	Symbol    string `json:"symbol"`    // 货币符号
	Multi     int64  `json:"multi"`     // 货币换算倍率
	Precision int32  `json:"precision"` // 0 表示整数，3 表示三位小数
}

var (
	confS = []*Conf{
		// 亚洲
		{"IDR", "Rp", 1, 0},
		{"PHP", "₱", multi, precision},
		{"INR", "₹", multi, precision},
		{"KRW", "₩", multi, precision},
		{"MYR", "RM", multi, precision},
		{"VND", "₫", 1, 0},
		{"MMK", "K", multi, precision},
		{"JPY", "¥", multi, precision},
		{"UZS", "UZS", 1, 0},
		{"KGS", "KGS", multi, precision},
		{"AZN", "₼", multi, precision},
		{"BDT", "৳", multi, precision},
		{"THB", "฿", multi, precision},
		{"KZT", "₸", multi, precision},
		{"KHR", "៛", multi, precision},
		{"PKR", "₨", multi, precision},
		{"AED", "AED", multi, precision},
		{"TRY", "₺", multi, precision},
		{"SGD", "S$", multi, precision},
		{"NPR", "₨", multi, precision},
		{"MNT", "₮", multi, precision},
		{"LKR", "Rs", multi, precision},
		{"TWD", "NT$", multi, precision},
		{"LAK", "₭", 1, 0},
		{"HKD", "HK$", multi, precision},
		{"AMD", "֏", multi, precision},
		{"BND", "B$", multi, precision},
		{"IRR", "IRR", 1, 0},
		{"SAR", "SR", multi, precision},

		// 欧洲
		{"EUR", "€", multi, precision},
		{"PLN", "zł", multi, precision},
		{"UAH", "₴", multi, precision},
		{"CHF", "CHF", multi, precision},
		{"SEK", "kr", multi, precision},
		{"GBP", "£", multi, precision},
		{"RUB", "₽", multi, precision},
		{"DKK", "kr", multi, precision},
		{"MDL", "MDL", multi, precision},
		{"NOK", "kr", multi, precision},
		{"CZK", "Kč", multi, precision},

		// 非洲
		{"ZAR", "R", multi, precision},
		{"XOF", "CFA", multi, precision},
		{"GNF", "GNF", multi, precision},
		{"RWF", "RWF", multi, precision},
		{"XAF", "FCFA", multi, precision},
		{"TZS", "TZS", multi, precision},
		{"EGP", "E£", multi, precision},
		{"GHS", "GH₵", multi, precision},
		{"UGX", "UGX", multi, precision},
		{"NGN", "₦", multi, precision},
		{"ZMW", "ZK", multi, precision},
		{"TND", "TND", multi, precision},
		{"MAD", "DH", multi, precision},
		{"DZD", "DA", multi, precision},
		{"ETB", "ETB", multi, precision},
		{"CDF", "CDF", multi, precision},
		{"KES", "KSh", multi, precision},
		{"NAD", "N$", multi, precision},

		// 美洲
		{"USD", "$", multi, precision},
		{"ARS", "AR$", multi, precision},
		{"CLP", "CLP$", multi, precision},
		{"COP", "COP$", multi, precision},
		{"MXN", "MX$", multi, precision},
		{"PEN", "S/.", multi, precision},
		{"CAD", "C$", multi, precision},
		{"BRL", "R$", multi, precision},
		{"BOB", "Bs.", multi, precision},
		{"HNL", "L", multi, precision},
		{"PYG", "₲", multi, precision},
		{"VES", "Bs.", multi, precision},

		// 大洋洲
		{"AUD", "A$", multi, precision},
		{"NZD", "NZ$", multi, precision},
		{"PGK", "PGK", multi, precision},

		// 其它
		{"GC", "GC", 1, 0}, // 虚拟币，不需小数
		{"SC", "SC", multi, precision},
		{TRC, "", 1, 0},
	}

	confM = map[string]*Conf{}
)

func init() {
	for _, item := range confS {
		confM[item.Currency] = item
	}
}

//goland:noinspection GoUnusedExportedFunction
func IsValidCurrency(currency string) bool {
	_, ok := confM[currency]
	return ok
}

//goland:noinspection GoUnusedExportedFunction
func GetConf(currency string) *Conf {
	return confM[currency]
}

//goland:noinspection GoUnusedExportedFunction
func GetConfS() []*Conf {
	return confS
}

//goland:noinspection GoUnusedExportedFunction
func GetConfM() map[string]*Conf {
	return confM
}
