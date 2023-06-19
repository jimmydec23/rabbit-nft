package util

import (
	"crypto/ecdsa"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

func GetPrivateKeyByHash(hash string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hash[2:])
}

func GetAddressByPrivate(private *ecdsa.PrivateKey) string {
	address := crypto.PubkeyToAddress(private.PublicKey).Hex()
	return address
}

func TokenId2Byte(tokenId uint64) []byte {
	return []byte(strconv.FormatUint(tokenId, 10))
}

func TokenIdFromByte(tokenIdByte []byte) uint64 {
	tokenId, _ := strconv.ParseUint(string(tokenIdByte), 10, 64)
	return tokenId
}

func TokenIdFromString(tokenIdStr string) uint64 {
	tokenId, _ := strconv.ParseUint(tokenIdStr, 10, 64)
	return tokenId
}
