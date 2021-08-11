package segment

import (
	"github.com/adamzy/cedar-go"
)

// Dictionary 结构体实现了一个字串前缀树，一个分词可能出现在叶子节点也有可能出现在非叶节点
type Dictionary struct {
	trie        *cedar.Cedar // 前缀树：Cedar 前缀树
	tokens      []Token      // 分词数组：词典中所有的分词，方便遍历
	maxTokenLen int          // 词典中最长的分词
	totalFreq   int64        // 词典中所有分词的频率之和
}

func NewDictionary() *Dictionary {
	return &Dictionary{trie: cedar.New()}
}

//// 词典中最长的分词
//func (dict *Dictionary) MaxTokenLength() int {
//	return dict.maxTokenLen
//}

// TokensLen 词典中分词数目
func (dict Dictionary) TokensLen() int {
	return len(dict.tokens)
}

// TotalFrequency 词典中所有分词的频率之和
func (dict Dictionary) TotalFrequency() int64 {
	return dict.totalFreq
}

// 向词典中加入一个分词
func (dict *Dictionary) addToken(token Token) {
	_, err := dict.trie.Get(token.text) // 已有分词
	if err == nil {
		return
	}

	_ = dict.trie.Insert(token.text, dict.TokensLen()) // 加入前缀树
	dict.tokens = append(dict.tokens, token)           // 加入分词数组

	// 统计
	dict.totalFreq += int64(token.freq)     // 总频率
	if len(token.text) > dict.maxTokenLen { // 最长分词
		dict.maxTokenLen = len(token.text)
	}
}

// 在词典中查找和字元组chars可以前缀匹配的所有分词
func (dict *Dictionary) lookupTokens(chars []Char) (tokens []*Token) {
	var id, value int // id = 0, from root
	var err error
	index := 0
	for _, char := range chars {
		// <id, value> + char = <id, value>
		id, err = dict.trie.Jump(char, id)
		if err != nil {
			break
		}
		value, err = dict.trie.Value(id)
		if err == nil {
			tokens = append(tokens, &dict.tokens[value])
			index++
		}
	}
	return
}
