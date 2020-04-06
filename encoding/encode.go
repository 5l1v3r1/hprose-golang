/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| encoding/encode.go                                       |
|                                                          |
| LastModified: Mar 21, 2020                               |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/modern-go/reflect2"
)

const (
	digits = "0123456789"
	digit2 = "" +
		"0001020304050607080910111213141516171819" +
		"2021222324252627282930313233343536373839" +
		"4041424344454647484950515253545556575859" +
		"6061626364656667686970717273747576777879" +
		"8081828384858687888990919293949596979899"
	digit3 = "" +
		"000001002003004005006007008009010011012013014015016017018019" +
		"020021022023024025026027028029030031032033034035036037038039" +
		"040041042043044045046047048049050051052053054055056057058059" +
		"060061062063064065066067068069070071072073074075076077078079" +
		"080081082083084085086087088089090091092093094095096097098099" +
		"100101102103104105106107108109110111112113114115116117118119" +
		"120121122123124125126127128129130131132133134135136137138139" +
		"140141142143144145146147148149150151152153154155156157158159" +
		"160161162163164165166167168169170171172173174175176177178179" +
		"180181182183184185186187188189190191192193194195196197198199" +
		"200201202203204205206207208209210211212213214215216217218219" +
		"220221222223224225226227228229230231232233234235236237238239" +
		"240241242243244245246247248249250251252253254255256257258259" +
		"260261262263264265266267268269270271272273274275276277278279" +
		"280281282283284285286287288289290291292293294295296297298299" +
		"300301302303304305306307308309310311312313314315316317318319" +
		"320321322323324325326327328329330331332333334335336337338339" +
		"340341342343344345346347348349350351352353354355356357358359" +
		"360361362363364365366367368369370371372373374375376377378379" +
		"380381382383384385386387388389390391392393394395396397398399" +
		"400401402403404405406407408409410411412413414415416417418419" +
		"420421422423424425426427428429430431432433434435436437438439" +
		"440441442443444445446447448449450451452453454455456457458459" +
		"460461462463464465466467468469470471472473474475476477478479" +
		"480481482483484485486487488489490491492493494495496497498499" +
		"500501502503504505506507508509510511512513514515516517518519" +
		"520521522523524525526527528529530531532533534535536537538539" +
		"540541542543544545546547548549550551552553554555556557558559" +
		"560561562563564565566567568569570571572573574575576577578579" +
		"580581582583584585586587588589590591592593594595596597598599" +
		"600601602603604605606607608609610611612613614615616617618619" +
		"620621622623624625626627628629630631632633634635636637638639" +
		"640641642643644645646647648649650651652653654655656657658659" +
		"660661662663664665666667668669670671672673674675676677678679" +
		"680681682683684685686687688689690691692693694695696697698699" +
		"700701702703704705706707708709710711712713714715716717718719" +
		"720721722723724725726727728729730731732733734735736737738739" +
		"740741742743744745746747748749750751752753754755756757758759" +
		"760761762763764765766767768769770771772773774775776777778779" +
		"780781782783784785786787788789790791792793794795796797798799" +
		"800801802803804805806807808809810811812813814815816817818819" +
		"820821822823824825826827828829830831832833834835836837838839" +
		"840841842843844845846847848849850851852853854855856857858859" +
		"860861862863864865866867868869870871872873874875876877878879" +
		"880881882883884885886887888889890891892893894895896897898899" +
		"900901902903904905906907908909910911912913914915916917918919" +
		"920921922923924925926927928929930931932933934935936937938939" +
		"940941942943944945946947948949950951952953954955956957958959" +
		"960961962963964965966967968969970971972973974975976977978979" +
		"980981982983984985986987988989990991992993994995996997998999"
)

var minInt64Buf = []byte("-9223372036854775808")

func toBytes(i uint64, buf []byte) (off int) {
	off = len(buf)
	var q, p uint64
	for i >= 100 {
		q = i / 1000
		p = (i - (q * 1000)) * 3
		i = q
		off -= 3
		copy(buf[off:off+3], digit3[p:p+3])
	}
	if i >= 10 {
		q = i / 100
		p = (i - (q * 100)) * 2
		i = q
		off -= 2
		copy(buf[off:off+2], digit2[p:p+2])
	}
	if i > 0 {
		off--
		buf[off] = digits[i]
	}
	return
}

// AppendInt64 i to buf
func AppendInt64(buf []byte, i int64) []byte {
	if i >= 0 {
		return AppendUint64(buf, uint64(i))
	}
	if i == math.MinInt64 {
		return append(buf, minInt64Buf...)
	}
	var u uint64 = uint64(-i)
	var buffer [20]byte
	off := toBytes(u, buffer[:]) - 1
	buffer[off] = '-'
	return append(buf, buffer[off:]...)
}

// AppendUint64 i to buf
func AppendUint64(buf []byte, i uint64) []byte {
	if (i >= 0) && (i <= 9) {
		return append(buf, digits[i])
	}
	var buffer [20]byte
	off := toBytes(i, buffer[:])
	return append(buf, buffer[off:]...)
}

// WriteInt64 to encoder
func WriteInt64(enc *Encoder, i int64) {
	if (i >= 0) && (i <= 9) {
		enc.buf = append(enc.buf, digits[i])
	} else {
		var tag = TagInteger
		if (i < math.MinInt32) || (i > math.MaxInt32) {
			tag = TagLong
		}
		enc.buf = append(enc.buf, tag)
		enc.buf = AppendInt64(enc.buf, i)
		enc.buf = append(enc.buf, TagSemicolon)
	}
}

// WriteUint64 to encoder
func WriteUint64(enc *Encoder, i uint64) {
	if (i >= 0) && (i <= 9) {
		enc.buf = append(enc.buf, digits[i])
	} else {
		var tag = TagInteger
		if i > math.MaxInt32 {
			tag = TagLong
		}
		enc.buf = append(enc.buf, tag)
		enc.buf = AppendUint64(enc.buf, i)
		enc.buf = append(enc.buf, TagSemicolon)
	}
}

// WriteInt32 to encoder
func WriteInt32(enc *Encoder, i int32) {
	if (i >= 0) && (i <= 9) {
		enc.buf = append(enc.buf, digits[i])
	} else {
		enc.buf = append(enc.buf, TagInteger)
		enc.buf = AppendInt64(enc.buf, int64(i))
		enc.buf = append(enc.buf, TagSemicolon)
	}
}

// WriteUint32 to encoder
func WriteUint32(enc *Encoder, i uint32) {
	WriteUint64(enc, uint64(i))
}

// WriteInt16 to encoder
func WriteInt16(enc *Encoder, i int16) {
	WriteInt32(enc, int32(i))
}

// WriteUint16 to encoder
func WriteUint16(enc *Encoder, i uint16) {
	if (i >= 0) && (i <= 9) {
		enc.buf = append(enc.buf, digits[i])
		return
	}
	enc.buf = append(enc.buf, TagInteger)
	enc.buf = AppendUint64(enc.buf, uint64(i))
	enc.buf = append(enc.buf, TagSemicolon)
	return
}

// WriteInt8 to encoder
func WriteInt8(enc *Encoder, i int8) {
	WriteInt32(enc, int32(i))
}

// WriteUint8 to encoder
func WriteUint8(enc *Encoder, i uint8) {
	WriteUint16(enc, uint16(i))
}

// WriteInt to encoder
func WriteInt(enc *Encoder, i int) {
	WriteInt64(enc, int64(i))
}

// WriteUint to encoder
func WriteUint(enc *Encoder, i uint) {
	WriteUint64(enc, uint64(i))
}

// WriteNil to encoder
func WriteNil(enc *Encoder) {
	enc.buf = append(enc.buf, TagNull)
}

// WriteBool to encoder
func WriteBool(enc *Encoder, b bool) {
	if b {
		enc.buf = append(enc.buf, TagTrue)
	} else {
		enc.buf = append(enc.buf, TagFalse)
	}
}

func writeFloat(enc *Encoder, f float64, bitSize int) {
	switch {
	case f != f:
		enc.buf = append(enc.buf, TagNaN)
	case f > math.MaxFloat64:
		enc.buf = append(enc.buf, TagInfinity, TagPos)
	case f < -math.MaxFloat64:
		enc.buf = append(enc.buf, TagInfinity, TagNeg)
	default:
		enc.buf = append(enc.buf, TagDouble)
		enc.buf = strconv.AppendFloat(enc.buf, f, 'g', -1, bitSize)
		enc.buf = append(enc.buf, TagSemicolon)
	}
}

// WriteFloat32 to encoder
func WriteFloat32(enc *Encoder, f float32) {
	writeFloat(enc, float64(f), 32)
}

// WriteFloat64 to encoder
func WriteFloat64(enc *Encoder, f float64) {
	writeFloat(enc, f, 64)
}

func utf16Length(str string) (n int) {
	length := len(str)
	n = length
	c := 0
	for i := 0; i < length; i++ {
		a := str[i]
		if c == 0 {
			switch {
			case (a & 0xe0) == 0xc0:
				c = 1
				n--
			case (a & 0xf0) == 0xe0:
				c = 2
				n -= 2
			case (a & 0xf8) == 0xf0:
				c = 3
				n -= 2
			case (a & 0x80) == 0x80:
				return -1
			}
		} else {
			if (a & 0xc0) != 0x80 {
				return -1
			}
			c--
		}
	}
	if c != 0 {
		return -1
	}
	return n
}

func appendBinary(buf []byte, bytes []byte, length int) []byte {
	if length > 0 {
		buf = AppendUint64(buf, uint64(length))
	}
	buf = append(buf, TagQuote)
	buf = append(buf, bytes...)
	buf = append(buf, TagQuote)
	return buf
}

func appendBytes(buf []byte, bytes []byte) []byte {
	buf = append(buf, TagBytes)
	buf = appendBinary(buf, bytes, len(bytes))
	return buf
}

func appendString(buf []byte, s string, length int) []byte {
	if length < 0 {
		return appendBytes(buf, reflect2.UnsafeCastString(s))
	}
	buf = append(buf, TagString)
	buf = appendBinary(buf, reflect2.UnsafeCastString(s), length)
	return buf
}

// WriteHead to encoder, n is the count of elements in list or map
func WriteHead(enc *Encoder, n int, tag byte) {
	enc.buf = append(enc.buf, tag)
	if n > 0 {
		enc.buf = AppendUint64(enc.buf, uint64(n))
	}
	enc.buf = append(enc.buf, TagOpenbrace)
}

// WriteObjectHead to encoder, r is the reference number of struct
func WriteObjectHead(enc *Encoder, r int) {
	enc.buf = append(enc.buf, TagObject)
	enc.buf = AppendUint64(enc.buf, uint64(r))
	enc.buf = append(enc.buf, TagOpenbrace)
}

// WriteFoot of list or map to encoder
func WriteFoot(enc *Encoder) {
	enc.buf = append(enc.buf, TagClosebrace)
}

func writeComplex(enc *Encoder, r float64, i float64, bitSize int) {
	if i == 0 {
		writeFloat(enc, r, bitSize)
	} else {
		enc.AddReferenceCount(1)
		WriteHead(enc, 2, TagList)
		writeFloat(enc, r, bitSize)
		writeFloat(enc, i, bitSize)
		WriteFoot(enc)
	}
}

// WriteComplex64 to encoder
func WriteComplex64(enc *Encoder, c complex64) {
	writeComplex(enc, float64(real(c)), float64(imag(c)), 32)
}

// WriteComplex128 to encoder
func WriteComplex128(enc *Encoder, c complex128) {
	writeComplex(enc, real(c), imag(c), 64)
}

// WriteBigInt to encoder
func WriteBigInt(enc *Encoder, i *big.Int) {
	enc.buf = append(enc.buf, TagLong)
	enc.buf = append(enc.buf, i.String()...)
	enc.buf = append(enc.buf, TagSemicolon)
}

// WriteBigFloat to encoder
func WriteBigFloat(enc *Encoder, f *big.Float) {
	enc.buf = append(enc.buf, TagDouble)
	enc.buf = f.Append(enc.buf, 'g', -1)
	enc.buf = append(enc.buf, TagSemicolon)
}

// WriteBigRat to encoder
func WriteBigRat(enc *Encoder, r *big.Rat) {
	if r.IsInt() {
		WriteBigInt(enc, r.Num())
	} else {
		enc.AddReferenceCount(1)
		s := r.String()
		enc.buf = appendString(enc.buf, s, len(s))
	}
}

// WriteError to encoder
func WriteError(enc *Encoder, e error) {
	enc.buf = append(enc.buf, TagError)
	enc.AddReferenceCount(1)
	s := e.Error()
	enc.buf = appendString(enc.buf, s, utf16Length(s))
}

// EncodeReference to encoder
func EncodeReference(valenc ValueEncoder, enc *Encoder, v interface{}) {
	if reflect2.IsNil(v) {
		WriteNil(enc)
	} else if ok := enc.WriteReference(v); !ok {
		valenc.Write(enc, v)
	}
}

// SetReference to encoder
func SetReference(enc *Encoder, v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		enc.SetReference(v)
	} else {
		enc.AddReferenceCount(1)
	}
}