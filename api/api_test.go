package api

import (
	"testing"
)

func TestCheckBVComplete(t *testing.T) {
	un := "1661628099@qq.com"
	speedID := "1HrwVZLRTGiK9ZsKmG8OqsBQdpAr1rCkE9zxOIIOvAgA"
	c, err := checkBVComplete(un, speedID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(c)
}
