package zerologutil

// 日志关键字数组
// 用于支持搜索查询的索引建立
// 以 [k1, v1, k2, v2 ...] 的形式 组装信息
// 使用时 会被转换成 key-value 组成的数组 来使用;
// key 必须为 string, 否则对应的一组 key-value 会被忽略
// kvList%2 != 0 时, 多余的 key 会被忽略

type FieldArr = []any

// Fields , FieldArr 的构建打包函数
//
//goland:noinspection GoUnusedExportedFunction
func Fields(kvList ...any) FieldArr {
	return kvList
}

// FieldsByMap map转FieldArr 的构建打包函数
//
//goland:noinspection GoUnusedExportedFunction
func FieldsByMap[K comparable](m map[K]any) FieldArr {
	var ret FieldArr
	for k, v := range m {
		ret = append(ret, k)
		ret = append(ret, v)
	}
	return ret
}
