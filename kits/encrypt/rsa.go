package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"os"
	"path"
)

type rsa_encrper struct {
	publickey  *rsa.PublicKey
	privatekey *rsa.PrivateKey
	hsh        hash.Hash
}

func (s *rsa_encrper) Sign(src []byte) ([]byte, error) {
	s.hsh.Reset()
	_, err := s.hsh.Write(src)
	var sic []byte
	if err == nil {
		bys := s.hsh.Sum(nil)
		sic, err = rsa.SignPSS(rand.Reader, s.privatekey, crypto.SHA256, bys, nil)
	}
	return sic, err

}
func (s *rsa_encrper) Verify(src, sig []byte) bool {
	s.hsh.Reset()
	_, err := s.hsh.Write(src)
	if err == nil {
		hshval := s.hsh.Sum(nil)
		err = rsa.VerifyPSS(s.publickey, crypto.SHA256, hshval, sig, nil)
	}
	return err == nil

}
func (s *rsa_encrper) Encrpyt(src []byte) (ans []byte, err error) {
	if s.publickey != nil {
		ans, err = rsa.EncryptPKCS1v15(rand.Reader, s.publickey, src)
	} else {
		err = fmt.Errorf("not set public key")
	}
	return
}

func (s *rsa_encrper) Decrpyt(src []byte) (ans []byte, err error) {
	if s.privatekey != nil {
		ans, err = rsa.DecryptPKCS1v15(rand.Reader, s.privatekey, src)
	} else {
		err = fmt.Errorf("not set private key")
	}
	return
}
func new_rsa_encrper(key ...[]byte) (Cryptor, error) {
	if len(key) < 1 {
		return nil, Lack_Args_ERR
	}
	block, _ := pem.Decode(key[0])
	bi, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		pk, ok := bi.(*rsa.PublicKey)
		if ok {
			return &rsa_encrper{publickey: pk}, nil
		}
	}
	return nil, err
}
func new_rsa_decrper(key ...[]byte) (Cryptor, error) {
	if len(key) < 1 {
		return nil, Lack_Args_ERR
	}
	block, _ := pem.Decode(key[0])
	bi, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return &rsa_encrper{privatekey: bi}, nil
	} else {
		return nil, err
	}
}

// [private,public]
func new_rsa_Cryptor(key ...[]byte) (Cryptor, error) {
	if len(key) < 2 {
		return nil, Lack_Args_ERR
	}
	var (
		ans           rsa_encrper
		err           error
		block_private *pem.Block
	)
	ans.hsh = sha256.New()
	block_private, _ = pem.Decode(key[0])
	if block_private == nil {
		err = fmt.Errorf("get block private pem failed")
		return nil, err
	}

	ans.privatekey, err = x509.ParsePKCS1PrivateKey(block_private.Bytes)
	if err == nil {
		block_public, _ := pem.Decode(key[1])
		if block_public == nil {
			err = fmt.Errorf("get block private pem failed")
			return nil, err
		}
		ans.publickey, err = x509.ParsePKCS1PublicKey(block_public.Bytes)
		// var newkey any
		// newkey, err = x509.ParsePKIXPublicKey(block_public.Bytes)
		// if err == nil {
		// 	// var ok bool
		// 	ans.publickey, ok = newkey.(*rsa.PublicKey)
		// 	if !ok {
		// 		err = fmt.Errorf("transfer key to public key failed")
		// 	}
		// }
	}
	if ans.publickey == nil {
		fmt.Println("[warn] public key is null")
	}
	if ans.privatekey == nil {
		fmt.Println("[warn] private key is null")
	}
	return &ans, err
}

func RSA_Sign(cipher Cryptor, src []byte) (ans []byte, err error) {
	if cipher == nil {
		err = fmt.Errorf("cipher is null")
		return
	}
	rsa_cipher, ok := cipher.(*rsa_encrper)
	if ok {
		ans, err = rsa_cipher.Sign(src)
	} else {
		err = fmt.Errorf("it's not rsa cipher")
	}
	return
}
func RSA_Verify(cipher Cryptor, src, sig []byte) error {
	if cipher == nil {
		return fmt.Errorf("cipher is nil")
	}
	rsa_cipher, ok := cipher.(*rsa_encrper)
	var err error
	if ok {
		if !rsa_cipher.Verify(src, sig) {
			err = fmt.Errorf("verify signature failed")
		}
	} else {
		err = fmt.Errorf("cipher is not rsa cipher")
	}
	return err
}

// generate rsa key pair
func Generate_RSA_KeyPair(parent_path string, bits int) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	var (
		errs               error
		fi                 *os.File
		fis                *os.File
		privatekey_content []byte
		publickey_content  []byte
	)
	if err == nil {
		pubk := key.PublicKey
		puby := x509.MarshalPKCS1PublicKey(&pubk)
		pubpm := pem.Block{Type: "RSA Public Key", Bytes: puby}
		bys := x509.MarshalPKCS1PrivateKey(key)
		pm := pem.Block{Type: "RSA Private Key", Bytes: bys}
		fi, err = os.OpenFile(path.Join(parent_path, "rsa_private_key.pem"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0640)
		fis, errs = os.OpenFile(path.Join(parent_path, "rsa_public_key.pem"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0640)
		if err == nil && errs == nil {
			err = pem.Encode(fi, &pm)
			errs = pem.Encode(fis, &pubpm)
			if err == nil {
				privatekey_content = pem.EncodeToMemory(&pm)
				if privatekey_content == nil {
					err = fmt.Errorf("private key out to memory failed")
				}
			}
			if errs == nil {
				publickey_content = pem.EncodeToMemory(&pubpm)
				if publickey_content == nil {
					errs = fmt.Errorf("public key out to memory failed")
				}
			}
		}
		fi.Close()
		fis.Close()
	}
	if err != nil && errs != nil {
		err = errors.Join(err, errs)
	} else if errs != nil {
		err = errs
	}
	return privatekey_content, publickey_content, err

}
