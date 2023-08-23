package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

var register_encrper = map[string]func(keypair ...any) (Cryptor, error){"des": new_des_Cryptor, "rsa": new_rsa_Cryptor, "rsa_en": new_rsa_encrper, "rsa_de": new_rsa_decrper}
var Lack_Args_ERR = fmt.Errorf("lack args")

type regfunchandle func(keypair ...any) (Cryptor, error)
type Cryptor interface {
	Encrpyt(src []byte) (ans []byte, err error)
	Decrpyt(src []byte) (ans []byte, err error)
}

func GetCryptor(name string, keypair ...any) (ans Cryptor, err error) {
	if regfunc, ok := register_encrper[name]; ok {
		ans, err = regfunc(keypair...)
		return
	}
	err = fmt.Errorf("not %s Cryptor", name)
	return nil, err
}

type des_Cryptor struct {
	block cipher.Block
	key   []byte
}

// 明文补码
func (s *des_Cryptor) pkc5padding(origin *[]byte) {
	padding := len(*origin) % s.block.BlockSize()
	padcontent := bytes.Repeat([]byte{byte(padding)}, padding)
	*origin = append(*origin, padcontent...)
}
func (s *des_Cryptor) pkc5(origin *[]byte) {
	padding := int((*origin)[len(*origin)-1])
	*origin = (*origin)[:len(*origin)-padding]
}
func (s *des_Cryptor) Encrpyt(src []byte) (ans []byte, err error) {
	s.pkc5padding(&src)
	block_mod := cipher.NewCBCEncrypter(s.block, s.key)
	ans = make([]byte, len(src))
	block_mod.CryptBlocks(ans, src)
	return
}

func (s *des_Cryptor) Decrpyt(src []byte) (ans []byte, err error) {
	block_mod := cipher.NewCBCDecrypter(s.block, s.key)
	ans = make([]byte, len(src))
	block_mod.CryptBlocks(ans, src)
	s.pkc5(&ans)
	return
}
func new_des_Cryptor(key ...any) (Cryptor, error) {
	if len(key) < 1 {
		return nil, Lack_Args_ERR
	}
	var block cipher.Block
	var err error
	block, err = des.NewCipher(key[0].([]byte))
	if err == nil {
		return &des_Cryptor{block: block, key: key[0].([]byte)}, nil
	}
	return nil, err
}

const (
	Force = 90
	Smart = 91
)

func Register_Cipher(name string, regfunc regfunchandle, mod int8) bool {
	if _, ok := register_encrper[name]; !ok || mod == Force {
		register_encrper[name] = regfunc
		return true
	}
	return false
}
