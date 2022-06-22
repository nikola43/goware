package decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
)

var Key rsa.PrivateKey



func init() {
	Key = rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: FromBase10(""), // yes, yes change all of those
			E: 65537,
		},
		D: FromBase10(""),
		Primes: []*big.Int{
			FromBase10(""),
			FromBase10(""),
		},
	}
	Key.Precompute()
}

func main() {
	key, err := hex.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	aes_key, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &Key, key, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Key: %x\n", aes_key)
}
