package exn

// map 扩展工具函数

//goland:noinspection GoUnusedExportedFunction
func MapKeyS[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
