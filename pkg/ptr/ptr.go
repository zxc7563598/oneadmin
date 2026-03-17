package ptr

// Uint64 将 uint64 转换为 *uint64
func Uint64(v uint64) *uint64 {
	return &v
}

// Uint64 将 int 转换为 *int
func Int(v int) *int {
	return &v
}

// Uint64 将 string 转换为 *string
func String(v string) *string {
	return &v
}

// Uint64 将 bool 转换为 *bool
func Bool(v bool) *bool {
	return &v
}

// deref 安全地将指针解引用为值，如果指针为nil则返回零值
func Deref[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}
