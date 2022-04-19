package format

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
)

var keyPool = sync.Pool{
	New: func() interface{} {
		b := &bytes.Buffer{}
		b.Grow(16)
		return b
	},
}

type key struct{}

var Key key

func (k key) Gen(indexes ...string) string {
	buf := keyPool.Get().(*bytes.Buffer)
	buf.WriteString(strings.Join(indexes, ":"))

	t := buf.String()
	buf.Reset()
	keyPool.Put(buf)
	return t
}

// CacheLanguage cache
func (k key) CacheLanguage(tid int64) string {
	return k.Gen("lang", strconv.FormatInt(tid, 10))
}

// BoltLanguage bolt
func (k key) BoltLanguage(tid int64) []byte {
	return []byte(k.Gen(strconv.FormatInt(tid, 10)))
}
