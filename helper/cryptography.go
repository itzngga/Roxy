package helper

import (
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"github.com/spf13/viper"
	"golang.org/x/crypto/blowfish"
)

var src = rand.New(rand.NewSource(time.Now().UnixNano()))
var XTCrypto *Cryptography

type Cryptography struct {
	blowfishSalt *blowfish.Cipher
}

func NewCryptography() *Cryptography {
	if XTCrypto != nil {
		return XTCrypto
	}
	salted, err := blowfish.NewCipher([]byte(viper.GetString("SECRET_KEY")))
	if err != nil {
		panic(err.Error())
	}
	XTCrypto = &Cryptography{blowfishSalt: salted}
	return XTCrypto
}

func (p *Cryptography) MakeEncrypt(value string) (string, error) {
	val, _ := base64.StdEncoding.DecodeString(value)
	mode := ecb.NewECBEncrypter(p.blowfishSalt)
	pt := make([]byte, len(val))
	mode.CryptBlocks(pt, val)
	padder := padding.NewPkcs5Padding()
	pt, err := padder.Unpad(pt)
	return base64.StdEncoding.EncodeToString(pt), err
}

func (p *Cryptography) MakeDecrypt(value string) (string, error) {
	val, _ := base64.StdEncoding.DecodeString(value)
	mode := ecb.NewECBDecrypter(p.blowfishSalt)
	pt := make([]byte, len(val))
	mode.CryptBlocks(pt, val)
	padder := padding.NewPkcs5Padding()
	pt, err := padder.Unpad(pt)
	return base64.StdEncoding.EncodeToString(pt), err
}

func GenerateRandomID(n int) string {
	b := make([]byte, (n+1)/2) // can be simplified to n/2 if n is always even

	if _, err := src.Read(b); err != nil {
		panic(err)
	}

	return strings.ToUpper("XT" + hex.EncodeToString(b)[:n])
}
