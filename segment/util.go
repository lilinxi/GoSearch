package segment

// 取两整数较小值
func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// 取两整数较大值
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 将英文词转化为小写
func toLower(text []byte) []byte {
	output := make([]byte, len(text))
	for i, t := range text {
		if t >= 'A' && t <= 'Z' {
			output[i] = t - 'A' + 'a'
		} else {
			output[i] = t
		}
	}
	return output
}
