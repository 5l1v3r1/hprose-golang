/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| encoding/list_decoder_test.go                            |
|                                                          |
| LastModified: Jun 27, 2020                               |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"container/list"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeList(t *testing.T) {
	sb := &strings.Builder{}
	enc := NewEncoder(sb).Simple(false)
	lst := list.New()
	lst.PushBack(1)
	lst.PushBack(2)
	lst.PushBack(3)
	enc.Encode(*lst)
	enc.Encode(lst)
	enc.Encode(nil)
	enc.Encode("")
	enc.Encode("hello")
	dec := NewDecoder(([]byte)(sb.String()))
	var l *list.List
	dec.Decode(&l)
	assert.Equal(t, lst, l)
	dec.Decode(&l)
	assert.Equal(t, lst, l)
	dec.Decode(&l)
	assert.Equal(t, (*list.List)(nil), l)
	dec.Decode(&l)
	assert.Equal(t, list.New(), l)
	dec.Decode(&l)
	assert.EqualError(t, dec.Error, `hprose/encoding: can not cast string to *list.List`)
}
