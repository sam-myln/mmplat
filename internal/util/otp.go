package util

import (
	"crypto/sha256"
	"encoding/base32"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func ValidateKeyPass(loginPass, passcode string) (bool, error) {
	key, err := KeyGen(loginPass)
	if err != nil {
		return false, err
	}
	return totp.Validate(passcode, key.Secret()), nil
}

const passes = 377

type MemoryFile struct {
	data string
}

func (fd *MemoryFile) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	n = copy(p, fd.data[:])
	return n, nil
}

func (fd *MemoryFile) Write(p []byte) (n int, err error) {
	fd.data = string(p)
	return len(p), nil
}

func (fd *MemoryFile) SetData(hash string) {
	fd.data = hash
}

func (fd *MemoryFile) GetData() string {
	return fd.data
}

//KeyGen login pass format user:pass
func KeyGen(LoginPass string) (*otp.Key, error) {
	tf := &MemoryFile{LoginPass}
	for i := 0; i < passes; i++ {
		hasher := sha256.New()
		hasher.Write([]byte(tf.GetData()))
		hash := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hasher.Sum(nil))
		tf.SetData(hash)
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "example.com",
		AccountName: LoginPass,
		Secret:      []byte(tf.GetData()[:20]),// len 20
	})
	if err != nil {
		return key, err
	}

	return key, nil
}
