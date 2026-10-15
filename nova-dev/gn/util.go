package gn

// 放一些没有具体分类的工具函数

func ArrFirst[T any](arr []T) T {
	var zero T
	if len(arr) > 0 {
		return arr[0]
	}
	return zero
}

func ArrFirstOrDefault[T any](arr []T, def T) T {
	if len(arr) > 0 {
		return arr[0]
	}
	return def
}

// SingleToSlice 零值等于没传
//
//goland:noinspection GoUnusedExportedFunction
func SingleToSlice[T comparable](s T) []T {
	var z T
	if s == z {
		return nil
	}

	return []T{s}
}

//goland:noinspection GoUnusedExportedFunction
func ShallowCopy[T any](v T) *T {
	return &v
}
