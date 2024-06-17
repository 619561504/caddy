/**
 * Created by wuhanjie on 2024/4/16 16:49
 */

package filter

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type MiddlewareTest struct {
	ContentType string `json:"content_type"`
	// Regex to specify which pattern to look up
	SearchPattern string `json:"search_pattern"`
	// A byte-array specifying the string used to replace matches
	Replacement []byte `json:"replacement"`

	MaxSize int    `json:"max_size"`
	Path    string `json:"path"`
}

func TestEncode(t *testing.T) {
	v := MiddlewareTest{
		ContentType:   "content-type",
		SearchPattern: "dd",
		Replacement:   []byte("dd"),
		MaxSize:       0,
		Path:          "dd",
	}
	ret, _ := json.Marshal(v)
	fmt.Printf("%s\n", string(ret))
}

func Test(t *testing.T) {
	str := "{\"content_type\": \"text/html\", \"search_pattern\": \"</body>\", \"replacement\": \"</body><script>alert(1)</script>\"}"

	dec := json.NewDecoder(strings.NewReader(str))

	var v MiddlewareTest

	dec.DisallowUnknownFields()
	err := dec.Decode(&v)
	fmt.Printf("err: %v\n", err)
}
