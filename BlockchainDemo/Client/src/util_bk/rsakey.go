package util_bk

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	//	"strings"
)

var (
	errPublicModulus       = errors.New("crypto/rsa: missing public modulus")
	errPublicExponentSmall = errors.New("crypto/rsa: public exponent too small")
	errPublicExponentLarge = errors.New("crypto/rsa: public exponent too large")
)

func encrypt_priv(priv *rsa.PrivateKey, c *big.Int) *big.Int {
	m := new(big.Int).Exp(c, priv.D, priv.N)
	return m
}

func checkPub(pub *rsa.PublicKey) error {
	if pub.N == nil {
		return errPublicModulus
	}
	if pub.E < 2 {
		return errPublicExponentSmall
	}
	if pub.E > 1<<31-1 {
		return errPublicExponentLarge
	}
	return nil
}

func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

func TestPublicPrivateKey() {
	key1 := 2048
	key2 := 2048
	prefix1 := "key1"
	prefix2 := "key2"

	gen(key1, "key1")
	gen(key2, "key2")
	priv := GetPrivateKey(prefix1 + "priv" + strconv.Itoa(key1) + ".txt")
	priv2 := GetPrivateKey(prefix2 + "priv" + strconv.Itoa(key2) + ".txt")
	publicFile := prefix1 + "pub" + strconv.Itoa(key1) + ".txt"
	//getRSA_PKIXPublicKey(publicFile)

	//pub, err := x509.ParsePKCS1PublicKey(pub_block.Bytes)
	//pub, err := x509.ParsePKIXPublicKey(pub_block.Bytes)

	rsaPub := GetRSA_PKIXPublicKey(publicFile)

	plain := "test plain aaa"
	fmt.Println("-->Encrypted by public key ")
	encPub, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(plain))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("enc_pub: " + hex.Dump(encPub))

	decPriv, err := rsa.DecryptPKCS1v15(nil, priv, encPub)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("dec_priv: " + hex.Dump(decPriv))

	fmt.Println("Sign test")
	message := []byte("message to be signed too long llllllllllllllllllllllllllll 1111111111111111111111")

	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	rng := rand.Reader
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rng, priv, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}

	fmt.Printf("Signature: %x\n", signature)
	//rsaPub
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
	} else {
		fmt.Println("Sign is valid")
	}
	//priv2
	signature2, err := rsa.SignPKCS1v15(rng, priv2, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}

	fmt.Printf("Signature2: %x\n", signature2)
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], signature2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
	} else {
		fmt.Println("Sign is valid")
	}

}
func readKeyFile(filePath string) *pem.Block {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	block, _ := pem.Decode(buf)

	return block
}
func GetPrivateKey(filePath string) *rsa.PrivateKey {
	block := readKeyFile(filePath)

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return priv
}
func GetRSA_PKIXPublicKey(filePath string) *rsa.PublicKey {

	pub_block := readKeyFile(filePath)

	pub, err := x509.ParsePKIXPublicKey(pub_block.Bytes)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return pub.(*rsa.PublicKey)
}
func GetRSA_PKCS1PublicKey(filePath string) *rsa.PublicKey {
	pub_block := readKeyFile(filePath)
	pub, err := x509.ParsePKCS1PublicKey(pub_block.Bytes)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return pub
}
func gen(bit int, prefix string) {
	r := rand.Reader

	priv, err := rsa.GenerateKey(r, bit)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fp, err := os.Create(fmt.Sprintf(prefix+"priv%d.txt", bit))
	defer fp.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pem.Encode(fp, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	fp, err = os.Create(fmt.Sprintf(prefix+"pub%d.txt", bit))
	defer fp.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	//pubASN1 := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pem.Encode(fp, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
}
