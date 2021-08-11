package segment

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Text 字串。比如"中国有十三亿人口"
type Text []byte

// Char 字元。比如"中"又如"国", 英文的一个字元是一个词
type Char []byte

// Token 分词。比如"中国"又如"人口"
type Token struct {
	text  Text   // 分词的字串
	chars []Char // 分词的字串，这实际上是个字元数组
	freq  int    // 分词在语料库中的词频
	pos   string // 词性标注

	// log2(总词频/该分词词频)，这相当于log2(1/p(分词))，用作动态规划中
	// 该分词的路径长度。求解prod(p(分词))的最大值相当于求解
	// sum(distance(分词))的最小值，这就是“最短路径”的来历。
	distance float32

	// 该分词文本的进一步分词划分，见Segments函数注释。
	segments []*Segment
}

func NewToken(text Text, freq int, pos string) Token {
	return Token{
		text:     text,
		chars:    text.splitTextToChars(),
		freq:     freq,
		pos:      pos,
		distance: 0,
		segments: nil,
	}
}

// 将字串划分成字元
func (text Text) splitTextToChars() []Char {
	output := make([]Char, 0, len(text)/3) // 中文三个字节一个字元，英文平均一般不超过，特殊情况全是单字符

	current := 0 // 当前索引

	// 识别出拉丁字母或数字串
	inAlphanumeric := true
	alphanumericStart := 0

	for current < len(text) {
		r, size := utf8.DecodeRune(text[current:]) // 识别出一个 utf-8 字符

		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) { // 当前是拉丁字母或数字（非中日韩文字）
			if !inAlphanumeric { // 开始：拉丁字母或数字串
				alphanumericStart = current
				inAlphanumeric = true
			}
		} else {
			if inAlphanumeric { // 结束：拉丁字母或数字串
				inAlphanumeric = false
				if current != 0 {
					output = append(output, toLower(text[alphanumericStart:current])) // 添加：拉丁字母或数字串
				}
			}
			output = append(output, []byte(text[current:current+size])) // 添加：中日韩文字串
		}

		current += size // 步进索引
	}

	// 处理最后一个字元是拉丁字母或数字的情况
	if inAlphanumeric {
		if current != 0 {
			output = append(output, toLower(text[alphanumericStart:current]))
		}
	}

	return output
}

func (token Token) String() string {
	return fmt.Sprintf("%s", string(token.text))
}
