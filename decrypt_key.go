package main

import (
	"crypto/rsa"
	"math/big"
)

func init() {
	Key = rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: fromBase10(""), // yes, yes change all of those
			E: 65537,
		},
		D: fromBase10(""),
		Primes: []*big.Int{
			fromBase10(""),
			fromBase10(""),
		},
	}
	Key.Precompute()
}

/*
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

*/
