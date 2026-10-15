package exn

// slice 扩展工具函数

//goland:noinspection GoUnusedExportedFunction
func SliceContains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
