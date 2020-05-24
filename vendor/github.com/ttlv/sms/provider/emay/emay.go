package emay

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/tidwall/gjson"
	"github.com/ttacon/libphonenumber"
	"github.com/ttlv/sms"
	"github.com/ttlv/sms/config"
)

type ecbEncrypter ecb

type EmayProvider struct {
	DB *gorm.DB
}

type Message struct {
	Mobile             string `json:"mobile"`
	Content            string `json:"content"`
	RequestTime        int64  `json:"requestTime"`
	RequestValidPeriod int    `json:"requestValidPeriod"`
}

func New() EmayProvider {
	cfg := config.MustGetConfig()
	db, err := gorm.Open("mysql", cfg.DB)
	if err != nil {
		panic(err)
	}
	return EmayProvider{DB: db}
}

func (provider EmayProvider) GetCode() string {
	return "Emay"
}

func (provider EmayProvider) AvailableCountries() []string {
	return []string{"CN"}
}

func (provider EmayProvider) Available(s *sms.SendParams) bool {
	b := sms.SmsBrand{}
	provider.DB.Where("name = ?", s.Brand).First(&b)
	if !b.EnableEmay || b.EmayAppKey == "" || b.EmayAppID == "" {
		return false
	}
	return true
}

func (provider EmayProvider) Send(params sms.SendParams) (string, string, error) {
	brand := sms.SmsBrand{}
	phonenumber, _ := libphonenumber.Parse(params.Phone, params.Country)
	params.Phone = strings.Replace(libphonenumber.Format(phonenumber, libphonenumber.E164), " ", "", -1)
	provider.DB.Where("name = ?", params.Brand).First(&brand)
	data, _ := json.Marshal(Message{
		Mobile:             params.Phone,
		Content:            params.Content,
		RequestTime:        time.Now().UnixNano() / 1000000,
		RequestValidPeriod: 30,
	})
	encryptData := AesEncrypt(string(data), brand.EmayAppKey)
	req, _ := http.NewRequest("POST", "http://emay.ocx.com/inter/sendSingleSMS", bytes.NewReader(encryptData))
	req.Header.Add("appId", brand.EmayAppID)
	client := &http.Client{Timeout: time.Duration(3 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	result := resp.Header.Get("result")
	if result != "SUCCESS" {
		return "", "", fmt.Errorf("error")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	decryptJsonData := AesDecrypt(body, []byte(brand.EmayAppKey))
	return string(decryptJsonData), gjson.Get(string(decryptJsonData), "smsId").String(), err
}

func AesDecrypt(crypted, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("err is:", err)
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData
}

func AesEncrypt(src, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return crypted
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return []byte{}
	}
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

//type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
