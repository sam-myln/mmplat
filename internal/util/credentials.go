package util

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

//Credentials login:pass
type Credentials map[string]string

type CredMngr struct {
	files []string
	cred  Credentials
}

func CreateCredMngr(sources ...string) *CredMngr {
	return &CredMngr{sources, make(Credentials)}
}

// UniqSign generated unique singature for user pass combination for Totp
func (mngr *CredMngr) UniqSign(fn func(user, pass string)) {
	for i, v := range mngr.cred {
		fn(i, v)
	}
}

func (mngr *CredMngr) Exists(user, pass string) bool {
	return mngr.cred[user] != "" && Equals(mngr.cred[user], pass)
}

func (mngr *CredMngr) ParseCredentials() {
	for _, file := range mngr.files {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if NotEmpty(line) {
				t := strings.Split(line, " ")
				if Empty(t) {
					panic(errors.New("credentials: parsing error: wrong file format"))
				}
				mngr.cred[t[0]] = t[1]
			}
		}
	}
}
