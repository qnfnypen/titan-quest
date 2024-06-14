package api

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerify(t *testing.T) {
	sign := "9dbad6470b5c2473d3b13e4db994ad3913613c12ec51c0a2acacaf6962f1f8186b02fbb5b4340da6dc241c2c87772d898f221d4102715da5d0a5cdacfbbf4271"
	nonce := "TitanNetWork(223217)"
	pk := "0385134dcf502df55620199fb1e4cda253e73e8e869396437d00158eb12505b98a"
	address := "titan185he9e38325regattm2lc7l2amflahmy67s8ew"

	bytePubKey, err := hex.DecodeString(pk)
	require.Nil(t, err)

	byteSignature, err := hex.DecodeString(sign)
	require.Nil(t, err)

	pubKey := secp256k1.PubKey{Key: bytePubKey}
	success, err := VerifyArbitraryMsg(address, nonce, byteSignature, pubKey)
	require.Nil(t, err)

	assert.True(t, success)
}
