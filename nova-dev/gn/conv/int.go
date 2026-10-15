package conv

import "math"

// FloatToInt 四舍五入
func FloatToInt[T Float](f T) int {
	if f >= 0 {
		return int(f + 0.5)
	}
	return int(f - 0.5)
}

//goland:noinspection GoUnusedExportedFunction
func IsEqualFloat[T Float](l, r T) bool {
	return math.Abs(float64(l)-float64(r)) < 0.000001
}
