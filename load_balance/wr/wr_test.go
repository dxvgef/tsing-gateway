package wr

import (
	"log"
	"strconv"
	"testing"
)

func TestNext(t *testing.T) {
	var (
		err        error
		obj        InstType
		upstreamID = "WR"
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
		if err = obj.Set(upstreamID, addr, weight); err != nil {
			t.Error(err.Error())
			return
		}
		log.Println(addr, weight)
	}

	for i := 0; i < 20; i++ {
		t.Log(obj.Next(upstreamID))
	}
}
