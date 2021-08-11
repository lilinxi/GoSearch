package segment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestText_splitTextToChars(t *testing.T) {
	var text Text

	text = []byte("中国有十三亿人口")
	assert.Equal(t, "中/国/有/十/三/亿/人/口/", charsToString(text.splitTextToChars()))

	text = []byte("GitHub is a web-based hosting service, for software development projects.")
	assert.Equal(t, "github/ /is/ /a/ /web/-/based/ /hosting/ /service/,/ /for/ /software/ /development/ /projects/./", charsToString(text.splitTextToChars()))

	text = []byte("is is is is is is")
	assert.Equal(t, "is/ /is/ /is/ /is/ /is/ /is/", charsToString(text.splitTextToChars()))

	text = []byte("中国雅虎Yahoo! China致力于，领先的公益民生门户网站。")
	assert.Equal(t, "中/国/雅/虎/yahoo/!/ /china/致/力/于/，/领/先/的/公/益/民/生/门/户/网/站/。/", charsToString(text.splitTextToChars()))
}

func charsToString(bytes []Char) (output string) {
	for _, b := range bytes {
		output += (string(b) + "/")
	}
	return
}
