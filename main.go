package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// We're using a 32 byte long secret key.
	// This is probably something you generate first
	// then put into and environment variable.
	secretKey string = "8e8c2771be5c2bb10d541a5bf6aa51203e0bce2d6d4fa267afd89a6e20df11f1"
)

type Plaintext struct {
	Plaintext string `json:"plaintext"`
}

type Ciphertext struct {
	Ciphertext string `json:"ciphertext"`
}

func encrypt(c *gin.Context) {
	var input Plaintext
	if c.BindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	key, err := hex.DecodeString(secretKey)
	go func() {
		runtime.LockOSThread()

		begin := time.Now()
		log.Print("Running 100% CPU for", len(input.Plaintext), "microseconds")
		for {
			// run 100%
			if time.Now().Sub(begin) > time.Duration(len(input.Plaintext))*time.Microsecond {
				break
			}
		}
		// sleep
		time.Sleep(time.Duration(len(input.Plaintext)) * time.Microsecond)
	}()

	aes, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(input.Plaintext), nil)

	beforenc := []byte(ciphertext)
	encoded := make([]byte, hex.EncodedLen(len(beforenc)))
	hex.Encode(encoded, beforenc)
	c.JSON(http.StatusOK, Ciphertext{Ciphertext: string(encoded)})
}

func decrypt(c *gin.Context) {
	var input Ciphertext
	if c.BindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Print(input)
	key, err := hex.DecodeString(secretKey)
	aes, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	log.Default().Print("before Nonce")
	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	decode, _ := hex.DecodeString(input.Ciphertext)
	nonce, ciphertext := decode[:nonceSize], decode[nonceSize:]

	log.Default().Print(nonce)
	log.Default().Print(ciphertext)
	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, Plaintext{Plaintext: string(plaintext)})
}

func main() {

	router := gin.Default()

	router.POST("/encrypt", encrypt)
	router.POST("/decrypt", decrypt)
	router.Run()
}
