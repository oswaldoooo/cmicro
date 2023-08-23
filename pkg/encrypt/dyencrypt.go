package encrypt

import (
	"encoding/base32"
	"math/rand"
	"strconv"

	"github.com/oswaldoooo/cmicro/kits/encrypt"
)

// base64(time stamp)+dyencrypt(content)
type dy_cipher struct {
	disnum int
}

func init() {
	encrypt.Register_Cipher("dy_cipher", new_dycipher, encrypt.Smart)
}
func new_dycipher(keys ...any) (ans encrypt.Cryptor, err error) {
	disnn := 100
	if len(keys) > 0 {
		disn, ok := keys[0].(int)
		if ok {
			disnn = disn
		}
	}
	ans = &dy_cipher{disnum: disnn}
	return
}
func (s *dy_cipher) Encrpyt(src []byte) (ans []byte, err error) {
	stamp := rand.Int()
	dis := stamp % s.disnum
	ans = make([]byte, 2*len(src))
	var total byte
	for key, ele := range src {
		total = ele + byte(dis)
		ans[2*key] = total / 2
		ans[2*key+1] = total - ans[key*2]
	}
	ans = append([]byte(base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(stamp)))), ans...)
	return ans, nil
}

func (s *dy_cipher) Decrpyt(src []byte) (ans []byte, err error) {
	var (
		stampstr      []byte
		stamp, disnum int
	)
	stampstr, err = base32.StdEncoding.DecodeString(string(src[:32]))
	if err == nil {
		stamp, err = strconv.Atoi(string(stampstr))
		if err == nil {
			disnum = stamp % s.disnum
			src = src[32:]
			ans = make([]byte, len(src)/2)
			for e := 0; e < len(src)/2; e++ {
				ans[e] = src[2*e] + src[2*e+1] - byte(disnum)
			}
		}
	}
	return
}
