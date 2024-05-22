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

func TestCheckRFComplete(t *testing.T) {
	un := "0x30c16b1c6e07b5f685ee668b9e69a28512f74159"
	speedID := "1pzDbFVmSppnvVW5uquVEVFhJd-fvZ6LIF5YqmeuaVlA"

	c, err := checkRFComplete(un, speedID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(c)
}
