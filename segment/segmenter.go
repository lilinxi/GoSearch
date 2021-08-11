package segment

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

// 分词器结构体
type SegmenterImpl struct {
	dict *Dictionary
}

const (
	minTokenFrequency = 2 // 仅从字典文件中读取大于等于此频率的分词
)

//// 该结构体用于记录Viterbi算法中某字元处的向前分词跳转信息
//type jumper struct {
//	minDistance float32
//	token       *Token
//}

//// 返回分词器使用的词典
//func (seg *SegmenterImpl) Dictionary() *Dictionary {
//	return seg.dict
//}

// LoadDictionary 从文件中载入词典，词典的格式为（每个分词一行）：分词文本 频率 词性
func (seg *SegmenterImpl) LoadDictionary(dict string) {
	// 1. 创建词典
	seg.dict = NewDictionary()
	// 2. 载入词典
	dictFile, err := os.Open(dict)
	defer dictFile.Close()
	if err != nil {
		log.Fatalf("无法载入字典文件：%s \n", dict)
	}
	// 3. 读取词典
	reader := bufio.NewReader(dictFile)
	var text Text  // 文本
	var freq int   // 频率
	var pos string // 词性，Part of speech

	// 逐行读入分词
	for {
		// Fscanf 从 reader 扫描文本，根据 format 参数指定的格式将成功读取的空白分隔的值保存进成功传递给本函数的参数。返回成功扫描的条目个数和遇到的任何错误。
		size, err := fmt.Fscanln(reader, &text, &freq, &pos)

		// 读取机制
		if size == 0 || err == io.EOF { // 文件结束
			break
		} else if err != nil { // 非 EOF 错误
			log.Fatalf("扫描分词错误：%v\n", err)
		} else if size < 2 { // 无效行
			continue
		} else if freq < minTokenFrequency { // 过滤频率太小的词
			continue
		}

		// 将分词添加到字典中
		token := NewToken(text, freq, pos)
		seg.dict.addToken(token)
	}

	// 4. 计算分词属性
	// 计算每个分词的路径值，路径值含义见Token结构体的注释
	logTotalFreq := float32(math.Log2(float64(seg.dict.totalFreq)))
	for _, token := range seg.dict.tokens {
		token.distance = logTotalFreq - float32(math.Log2(float64(token.freq)))
		//fmt.Printf("%+v, %+v, %+v, %+v, %+v\n", seg.dict.totalFreq, logTotalFreq, token.freq, math.Log2(float64(token.freq)), token) // TODO
	}

	// 对每个分词进行细致划分，用于搜索引擎模式，该模式用法见Token结构体的注释。
	for i := range seg.dict.tokens {
		token := &seg.dict.tokens[i]
		fmt.Printf("token: %v\n", token)
		segments := seg.segmentChars(token.chars, true)

		// 计算需要添加的子分词数目
		numTokensToAdd := 0
		for iToken := 0; iToken < len(segments); iToken++ {
			if len(segments[iToken].token.text) > 0 {
				numTokensToAdd++
			}
		}
		token.segments = make([]*Segment, numTokensToAdd)

		// 添加子分词
		iSegmentsToAdd := 0
		for iToken := 0; iToken < len(segments); iToken++ {
			if len(segments[iToken].token.text) > 0 {
				token.segments[iSegmentsToAdd] = &segments[iToken]
				iSegmentsToAdd++
			}
		}
	}

	log.Println("sego词典载入完毕")
}

//// 对文本分词
////
//// 输入参数：
////	bytes	UTF8文本的字节数组
////
//// 输出：
////	[]Segment	划分的分词
//func (seg *Segmenter) Segment(bytes []byte) []Segment {
//	return seg.internalSegment(bytes, false)
//}
//
//
//func (seg *Segmenter) InternalSegment(bytes []byte, searchMode bool) []Segment {
//	return seg.internalSegment(bytes, searchMode)
//}
//
//func (seg *Segmenter) internalSegment(bytes []byte, searchMode bool) []Segment {
//	// 处理特殊情况
//	if len(bytes) == 0 {
//		return []Segment{}
//	}
//
//	// 划分字元
//	text := splitTextToChars(bytes)
//
//	return seg.segmentChars(text, searchMode)
//}
//

// 该结构体用于记录Viterbi算法中某字元处的向前分词跳转信息
type jumper struct {
	minDistance float32
	token       *Token
}

func (seg *SegmenterImpl) segmentChars(text []Char, searchMode bool) []Segment {
	// 搜索模式下该分词已无继续划分可能的情况
	if searchMode && len(text) == 1 {
		return []Segment{}
	}

	// jumpers定义了每个字元处的向前跳转信息，包括这个跳转对应的分词，
	// 以及从文本段开始到该字元的最短路径值
	//jumpers := make([]jumper, len(text))

	fmt.Println("text", text)

	//tokens := make([]*Token, seg.dict.maxTokenLen)
	for current := 0; current < len(text); current++ {
		// 找到前一个字元处的最短路径，以便计算后续路径值
		//var baseDistance float32
		//if current == 0 {
		//	// 当本字元在文本首部时，基础距离应该是零
		//	baseDistance = 0
		//} else {
		//	baseDistance = jumpers[current-1].minDistance
		//}

		// 寻找所有以当前字元开头的分词
		tokens := seg.dict.lookupTokens(text[current:minInt(current+seg.dict.maxTokenLen, len(text))])
		//numTokens := len(tokens)
		fmt.Println("-----")
		for t := range tokens {
			fmt.Println(t)
		}
		os.Exit(1)

		//// 对所有可能的分词，更新分词结束字元处的跳转信息
		//for iToken := 0; iToken < numTokens; iToken++ {
		//	location := current + len(tokens[iToken].text) - 1
		//	if !searchMode || current != 0 || location != len(text)-1 {
		//		updateJumper(&jumpers[location], baseDistance, tokens[iToken])
		//	}
		//}
		//
		//// 当前字元没有对应分词时补加一个伪分词
		//if numTokens == 0 || len(tokens[0].text) > 1 {
		//	updateJumper(&jumpers[current], baseDistance,
		//		&Token{text: []Text{text[current]}, frequency: 1, distance: 32, pos: "x"})
		//}
	}

	//// 从后向前扫描第一遍得到需要添加的分词数目
	numSeg := 0
	//for index := len(text) - 1; index >= 0; {
	//	location := index - len(jumpers[index].token.text) + 1
	//	numSeg++
	//	index = location - 1
	//}
	//
	//// 从后向前扫描第二遍添加分词到最终结果
	outputSegments := make([]Segment, numSeg)
	//for index := len(text) - 1; index >= 0; {
	//	location := index - len(jumpers[index].token.text) + 1
	//	numSeg--
	//	outputSegments[numSeg].token = jumpers[index].token
	//	index = location - 1
	//}
	//
	//// 计算各个分词的字节位置
	//bytePosition := 0
	//for iSeg := 0; iSeg < len(outputSegments); iSeg++ {
	//	outputSegments[iSeg].start = bytePosition
	//	//bytePosition += textSliceByteLength(outputSegments[iSeg].token.text)
	//	outputSegments[iSeg].end = bytePosition
	//}
	return outputSegments
}

// 更新跳转信息:
// 	1. 当该位置从未被访问过时(jumper.minDistance为零的情况)，或者
//	2. 当该位置的当前最短路径大于新的最短路径时
// 将当前位置的最短路径值更新为baseDistance加上新分词的概率
func updateJumper(jumper *jumper, baseDistance float32, token *Token) {
	newDistance := baseDistance + token.distance
	if jumper.minDistance == 0 || jumper.minDistance > newDistance {
		jumper.minDistance = newDistance
		jumper.token = token
	}
}
