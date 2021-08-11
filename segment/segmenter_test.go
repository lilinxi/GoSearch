package segment

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSegmenterImpl_LoadDictionary(t *testing.T) {
	pwd, _ := os.Getwd()
	dict := filepath.Join(pwd, "dict", "dictionary.txt")
	var segmenter Segmenter
	segmenter = &SegmenterImpl{}
	segmenter.LoadDictionary(dict)
}
