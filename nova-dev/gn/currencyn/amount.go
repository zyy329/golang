package currencyn

import (
	"strconv"
	"strings"

	"nova/gn/conv"
)

/*
note:
   有小数的，位数太多会有精度丢失，整数位只能支持到13位
   底层不处理异常，不支持的货币在上层过滤
*/

func FromStr(amount string, currency string) (i int64, e error) {
	// 截取小数点后超出的
	partS := strings.Split(amount, ".")
	if len(partS) == 2 {
		if more := len(partS[1]) - int(confM[currency].Precision); more > 0 {
			amount = amount[:len(amount)-more]
		}
	}

	var f float64
	if f, e = strconv.ParseFloat(amount, 64); e == nil {
		i = int64(conv.FloatToInt(f * float64(confM[currency].Multi)))
	}
	return
}

func ToStr(amount int64, currency string) string {
	return strconv.FormatFloat(float64(amount)/float64(confM[currency].Multi), 'f', -1, 64)
}
