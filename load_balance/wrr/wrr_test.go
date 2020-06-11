package wrr

import (
	"log"
	"strconv"
	"testing"
)

func TestNext(t *testing.T) {
	var (
		obj        InstType
		upstreamID = "WRR"
		addr       string
		weight     int
	)
	for i := 1; i <= 10; i++ {
		if i == 6 {
			addr = "10.0.0.6"
			weight = 6
		} else {
			addr = "10.0.0." + strconv.Itoa(i)
			weight = 1
		}
		obj.Set(upstreamID, addr, weight)
		log.Println(addr, weight)
	}

	for i := 0; i < 10; i++ {
		t.Log(obj.Next(upstreamID))
	}
}
