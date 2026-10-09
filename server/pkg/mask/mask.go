package mask

// Phone 手机号脱敏：138****8000
func Phone(p string) string {
	r := []rune(p)
	if len(r) == 11 {
		return string(r[:3]) + "****" + string(r[7:])
	}
	return p
}
