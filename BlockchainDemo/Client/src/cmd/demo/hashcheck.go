package demo

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"strconv"
)

func CheckHash(inputStr string) {
	aStringToHash := []byte(inputStr)
	//For CRC
	crc32Table := crc32.MakeTable(0xD5828281)
	crc64Table := crc64.MakeTable(0xC96C5795D7870F42)
	//For FNV
	fnvHash := fnv.New32()

	//Get the hashes in bytes
	md5Bytes := md5.Sum(aStringToHash)
	sha1Bytes := sha1.Sum(aStringToHash)
	sha256Bytes := sha256.Sum256(aStringToHash)
	sha512Bytes := sha512.Sum512(aStringToHash)
	fnvBytes := fnvHash.Sum(aStringToHash)

	//Get hashes in integers
	adler32Int := adler32.Checksum(aStringToHash)
	crc32Int := crc32.Checksum(aStringToHash, crc32Table)
	crc64Int := crc64.Checksum(aStringToHash, crc64Table)

	//Print out what will be hashed
	fmt.Println(string(aStringToHash))

	//Bytes to string
	fmt.Println("MD5 String is ", hex.EncodeToString(md5Bytes[:]))
	fmt.Println("SHA1 String is ", hex.EncodeToString(sha1Bytes[:]))
	fmt.Println("SHA256 String is ", hex.EncodeToString(sha256Bytes[:]))
	fmt.Println("SHA512 String is ", hex.EncodeToString(sha512Bytes[:]))
	fmt.Println("FNV String is ", hex.EncodeToString(fnvBytes[:]))

	//Uint to string
	fmt.Println("Adler32 String is ", strconv.FormatUint(uint64(adler32Int), 16))
	fmt.Println("CRC32 String is", strconv.FormatUint(uint64(crc32Int), 16))
	fmt.Println("CRC64 String is ", strconv.FormatUint(crc64Int, 16))
}
