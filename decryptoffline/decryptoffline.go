package decryptoffline

import (
	"encoding/hex"
)

func DecodeKey(encKey string) *[32]byte {
	key := [32]byte{}
	dKey, err := hex.DecodeString(encKey)
	if err != nil {
		panic(err)
	}
	copy(key[:], dKey)
	return &key
}

var key string = "" // your aes key (in hex)

/*
func main() {
	var files []string
	var counter int = 1
	var home string

	decryptionKey := DecodeKey(key)

	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	err := filepath.Walk(home, visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Printf("\rDecrypting %d/%d: %s", counter, len(files), file)

		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		decrypted, err := Decrypt(data, decryptionKey)
		if err != nil {
			log.Println(err)
			continue
		}

		err = ioutil.WriteFile(file, decrypted, 0644)
		if err != nil {
			continue
		}
		counter++
	}
	fmt.Printf("\n%d files decrypted.\n", len(files))
}

*/
