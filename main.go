package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"math/big"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	start            = time.Now()
	secretKey string = "8e8c2771be5c2bb10d541a5bf6aa51203e0bce2d6d4fa267afd89a6e20df11f1"
	timestamp string
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

	go func() {
		runtime.LockOSThread()
		// Generate a random string of 320000 bytes
		_, _ = GenerateRandomString(320000)
		time.Sleep(1 * time.Second)

	}()

	key, err := hex.DecodeString(secretKey)
	if err != nil {
		panic(err)
	}
	log.Print(" key is %s", key)
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

	beforenc := append([]byte(timestamp), []byte(ciphertext)...)
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

	decode, _ := hex.DecodeString(input.Ciphertext)

	ts, encrypted := decode[:8], decode[8:]
	log.Print("ts is ", string(ts))
	if string(ts) != timestamp {
		log.Print("expired")
		c.JSON(http.StatusGone, Plaintext{Plaintext: "expired"})
		return
	}

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
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	log.Default().Print(nonce)
	log.Default().Print(ciphertext)
	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, Plaintext{Plaintext: string(plaintext)})
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return hex.EncodeToString(ret), nil
}

func main() {

	go func() {
		runtime.LockOSThread()
		for {
			now := time.Now()

			if now.Sub(start) > 30*time.Second {
				secretKey, _ = GenerateRandomString(32)
				timestamp = strconv.FormatInt(now.Unix(), 16)
				log.Print("Time passed, generating a new secret key ", secretKey)
				log.Print("timestamp is ", timestamp)
				start = now
			}
			time.Sleep(1 * time.Second)
		}
	}()

	router := gin.Default()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.POST("/encrypt", encrypt)
	router.POST("/decrypt", decrypt)
	router.Run()
}
