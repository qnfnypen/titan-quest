package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func VerifyCosmosSign(nonce, sign, pk string) (bool, error) {
	bs, err := hex.DecodeString(pk)
	if err != nil {
		return false, fmt.Errorf("decode public key from pk error:%w", err)
	}

	pubkey := &secp256k1.PubKey{Key: bs}
	return VerifySignature(nonce, sign, pubkey)
}

func VerifyCosmosAddr(paddr, pk, pr string) (bool, error) {
	bs, err := hex.DecodeString(pk)
	if err != nil {
		return false, fmt.Errorf("decode public key from pk error:%w", err)
	}

	pubkey := &secp256k1.PubKey{Key: bs}
	addr, err := sdk.Bech32ifyAddressBytes(pr, pubkey.Address())
	if err != nil {
		return false, fmt.Errorf("get address error:%w", err)
	}

	return strings.EqualFold(paddr, addr), nil
}

// VerifySignature verifies a signed message using the provided public key
func VerifySignature(message, signatureHex string, pubKey cryptotypes.PubKey) (bool, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}
	return pubKey.VerifySignature(hash[:], signature), nil
}
