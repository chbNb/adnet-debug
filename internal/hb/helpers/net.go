package helpers

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	Aws    = "aws"
	Aliyun = "aliyun"
)

func IsCorrectIp(ip string) bool {
	if ip := net.ParseIP(ip); ip != nil {
		return true
	}
	return false
}

func QueryServerIp(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetCloudName() string {
	err := reqAws()
	if err == nil {
		return Aws
	}
	err = reqAliyun()
	if err == nil {
		return Aliyun
	}
	return Aliyun
}

func reqAws() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	uri := "http://169.254.169.254/latest/meta-data/services/domain"
	// req, _ := http.NewRequestWithContext(ctx, "GET", uri, nil)
	req, _ := http.NewRequest("GET", uri, nil)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func reqAliyun() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	uri := "http://100.100.100.200/latest/meta-data/eipv4"
	// req, _ := http.NewRequestWithContext(ctx, "GET", uri, nil)
	req, _ := http.NewRequest("GET", uri, nil)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
