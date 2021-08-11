package main

import (
	"fmt"
	"github.com/huichen/sego"
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
	"log"
)

var (
	// searcher是协程安全的
	searcher = engine.Engine{}
)

func main() {
	// 初始化
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: "/Users/limengfan/go/pkg/mod/github.com/huichen/wukong@v0.0.0-20161011030038-d014a1f19dae/data/dictionary.txt"})
	defer searcher.Close()

	// 将文档加入索引，docId 从1开始
	searcher.IndexDocument(1, types.DocumentIndexData{Content: "此次百度收购将成中国互联网最大并购"}, false)
	searcher.IndexDocument(2, types.DocumentIndexData{Content: "百度宣布拟全资收购91无线业务"}, false)
	searcher.IndexDocument(3, types.DocumentIndexData{Content: "百度是中国最大的搜索引擎"}, false)

	// 等待索引刷新完毕
	searcher.FlushIndex()

	// 搜索输出格式见types.SearchResponse结构体
	log.Print(searcher.Search(types.SearchRequest{Text:"百度中国"}))

	///////////////////////////////////////////////////////////////////////////////////////

	// 载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("/Users/limengfan/go/pkg/mod/github.com/huichen/wukong@v0.0.0-20161011030038-d014a1f19dae/data/dictionary.txt")
	fmt.Printf("%+v\n",segmenter.Dictionary().TotalFrequency())

	// 分词
	text := []byte("中华人民共和国中央人民政府")
	segments := segmenter.Segment(text)

	fmt.Printf("%+v\n", segments)
	fmt.Printf("%+v\n", segments[0].Token())

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	fmt.Println(sego.SegmentsToString(segments, false))
	fmt.Println(sego.SegmentsToString(segments, true))

	///////////////////////////////////////////////////////////////////////////////////////////

	//// create a new cedar trie.
	//trie := cedar.New()
	//
	//// a helper function to print the id-key-value triple given trie node id
	//printIdKeyValue := func(id int) {
	//	// the key of node `id`.
	//	key, _ := trie.Key(id)
	//	// the value of node `id`.
	//	value, _ := trie.Value(id)
	//	fmt.Printf("%d\t%s:%v\n", id, key, value)
	//}
	//
	//// Insert key-value pairs.
	//// The order of insertion is not important.
	//trie.Insert([]byte("How many"), 0)
	//trie.Insert([]byte("How many loved"), 1)
	//trie.Insert([]byte("How many loved your moments"), 2)
	//trie.Insert([]byte("How many loved your moments of glad grace"), 3)
	//trie.Insert([]byte("姑苏"), 4)
	//trie.Insert([]byte("姑苏城外"), 5)
	//trie.Insert([]byte("姑苏城外寒山寺"), 6)
	//
	//// Get the associated value of a key directly.
	//value, err := trie.Get([]byte("How many loved your moments of glad grace"))
	//fmt.Println("value", value,err)
	//
	//value, err = trie.Get([]byte("How"))
	//fmt.Println("value", value,err)
	//
	//value, err = trie.Get([]byte("How man"))
	//fmt.Println("value", value,err)
	//
	//value, err = trie.Get([]byte("How xxx"))
	//fmt.Println("value", value,err)
	//
	//value, err = trie.Get([]byte("How many l"))
	//fmt.Println("value", value,err)
	//
	//value, err = trie.Get([]byte("How many lxxx"))
	//fmt.Println("value", value,err)
	//
	//value, err = trie.Get([]byte("How many loved"))
	//fmt.Println("value", value,err)
	//
	//
	//// Or, jump to the node first,
	//id, _ := trie.Jump([]byte("How many loved your moments"), 0)
	//// then get the key and the value
	//printIdKeyValue(id)
	//
	//fmt.Println("\nPrefixMatch\nid\tkey:value")
	//for _, id := range trie.PrefixMatch([]byte("How many loved your moments of glad grace"), 0) {
	//	printIdKeyValue(id)
	//}
	//
	//fmt.Println("\nPrefixPredict\nid\tkey:value")
	//for _, id := range trie.PrefixPredict([]byte("姑苏"), 0) {
	//	printIdKeyValue(id)
	//}

}
