package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

type Payload struct {
	Id        string
	EpochTime uint64
	Distance  []float32
	Velocity  []float64
	Mask      []byte
}

func newPayload() *Payload {
	p := Payload{}
	p.Distance = make([]float32, 0)
	p.Velocity = make([]float64, 0)
	p.Mask = make([]byte, 0)
	return &p
}

// UnmarshalJSON deserializes ByteArray to object
func (payload *Payload) UnmarshalJSON(input []byte) {

	// ts, _ := jsonparser.GetInt(input, "a", "[0]", "ts")
	// tm, _ := jsonparser.GetInt(input, "a", "[0]", "tm")
	// ms := time.Unix(ts, tm).UTC()
	// fmt.Printf("%v\n", ms)

	dStr, _ := jsonparser.GetUnsafeString(input, "a", "[0]", "d")
	dStr = strings.Replace(dStr, "[", "", 1)
	dStr = strings.Replace(dStr, "]", "", len(dStr))

	slice := strings.Split(dStr, ",")
	for _, str := range slice {
		// 585a0184
		s1 := str[4:]
		s2 := str[0:4]

		v, _ := strconv.ParseInt(s1, 16, 64)
		vt := float64(v) / math.Pow(2, 8)
		if vt > (math.MaxUint8 / 2) {
			vt = vt - math.Pow(2, 8)
		}
		payload.Velocity = append(payload.Velocity, vt)

		d, _ := strconv.ParseInt(s2, 16, 64)
		xtt := d & 0x3FFF
		xt := float32(xtt) * 0.1
		payload.Distance = append(payload.Distance, xt)

		by := (byte)(((d & 0xC000) & 0xC000) >> 14) // get 2 MSB bits from the short
		payload.Mask = append(payload.Mask, by)

		/* 585a0184 vt: 1.5156, xt: 623.4 */
		//fmt.Printf("data:%s, velocity:%v, distance:%v, mask:%v\n", str, vt, xt, by)
	}
}
