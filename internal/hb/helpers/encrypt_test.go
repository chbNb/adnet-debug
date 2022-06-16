package helpers

import (
	"fmt"
	"testing"
)

func TestNewAESECBEncrypt(t *testing.T) {
	ecb := NewAESECBEncrypt([]byte("wMOs8JrslMwnbK44"), 16)
	raw := "af_cost_currency=USD&af_cost_model=CPI&af_cost_value=1.00"
	encryptStr, err := ecb.Encrypt(raw)
	if err != nil {
		t.Errorf("Encrypt err:%v", err)
	}

	rawVal, err := ecb.Decrypt(encryptStr)
	if err != nil {
		t.Errorf("Decrypt  err:%v", err)
	}
	if rawVal != raw {
		t.Errorf("Decrypt value unequal raw")
	}
}

func TestNewAESCBCEncrypt(t *testing.T) {
	// consumerKey, _ := hex.DecodeString("15d09e2d839ae599992bcfd48a20e3d0")
	// privateKey, _ := hex.DecodeString("b4aa405de27c5c8dd0fea55d482ddcca")
	privateKey := "15d09e2d839ae599992bcfd48a20e3d0"[:16]
	consumerKey := "b4aa405de27c5c8dd0fea55d482ddcca"

	cbc := NewAESCBCEncrypt([]byte(consumerKey), []byte(privateKey))
	raw := "cost=2.118&cost_model=cpi"
	encryptStr, err := cbc.Encrypt(raw)
	if err != nil {
		t.Errorf("Encrypt err:%v", err)
	}
	fmt.Println(encryptStr)
	rawVal, err := cbc.Decrypt(encryptStr)
	if err != nil {
		t.Errorf("Decrypt  err:%v", err)
	}

	if rawVal != raw {
		t.Errorf("Decrypt value unequal raw")
	}
}
