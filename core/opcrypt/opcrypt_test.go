package opcrypt

import (
	"fmt"
	"testing"
	"time"
)

var (
	key = "c8c29addd01f4c0eaefe23cbf1ac4943"
)

func TestAesEncryptCBC(t *testing.T) {
	data := fmt.Sprintf("0x30c16b1c6e07b5f685ee668b9e69a28512f74159:%d", time.Now().Unix())
	str, err := AesEncryptCBC([]byte(data), []byte(key))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(str)

	da, err := AesDecryptCBC(str, []byte(key))
	if err != nil {
		t.Fatal(err)
	}

	if string(da) != data {
		t.Fail()
	}
}
