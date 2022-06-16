package utility

import (
	"crypto/tls"
	"errors"
	"math"
	"net"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// timeout 超过100时，尽量控制在整百的粒度，让连接尽量在一个池里面
type FastHttpClient struct {
	timeout     int
	keepalive   time.Duration
	idles       int
	clients     [1001]*fasthttp.Client
	once        [1001]sync.Once
	moreClients map[int]*fasthttp.Client
}

func NewFastHttpClient(timeout, keepalive, idle int) *FastHttpClient {
	return &FastHttpClient{
		timeout:   timeout,
		keepalive: time.Second * time.Duration(keepalive),
		idles:     idle,
	}
}

// 最多支持10000毫秒的响应
var INVALID_TIMEOUT_PARAM = errors.New("invalid timeout param")

func (this *FastHttpClient) checkTimeout(timeout int) error {
	if timeout <= 10 || timeout > 10000 {
		return INVALID_TIMEOUT_PARAM
	}
	return nil
}

func (this *FastHttpClient) initClient(index int) {
	var timeout = time.Millisecond * time.Duration(index*10)
	this.clients[index] = &fasthttp.Client{
		Name:      "MTG SSP NET",                         // 支持自定义ua
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, // 支持ssl
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialDualStackTimeout(addr, timeout)
		}, // 1.支持双拨号；2.支持dns缓存；3.支持ip遍历; 4.支持超时控制
		DialDualStack:       true,           // 开启tcp4,tcp6双拨号支持
		MaxConnsPerHost:     this.idles,     // 最大分host个数
		MaxConnDuration:     this.keepalive, // 最大保持连接时间
		MaxIdleConnDuration: this.keepalive, // 最大空连接驻留时间
		ReadTimeout:         timeout,        // 读超时
		WriteTimeout:        timeout,        // 写超时
		MaxResponseBodySize: 1 << 20,        // 最大写1M
	}
}

func (this *FastHttpClient) TimeoutClient(timeout int) (client *fasthttp.Client, err error) {
	if timeout == 0 {
		timeout = this.timeout
	}
	if timeout <= 10 || timeout > 10000 {
		err = INVALID_TIMEOUT_PARAM
		return
	}
	var index int
	// 200ms以内，以10位单位取整，超过则以100位单位取整
	//if timeout < 200 {
	index = int(math.Ceil(float64(timeout) / 10))
	//} else {
	//	index = int(math.Ceil(float64(timeout)/30)) * 3
	//}
	if this.clients[index] == nil {
		this.once[index].Do(func() { this.initClient(index) })
	}
	client = this.clients[index]
	return
}

func (this *FastHttpClient) Client() (*fasthttp.Client, error) {
	return this.TimeoutClient(this.timeout)
}

func (this *FastHttpClient) Do(timeout int, req *fasthttp.Request, res *fasthttp.Response) error {
	client, err := this.TimeoutClient(timeout)
	if err != nil {
		return err
	}
	return client.Do(req, res)
}

var defaultFastHttpClient *FastHttpClient

func InitFastHttpClient(timeout, keepalive, idles int) {
	defaultFastHttpClient = NewFastHttpClient(timeout, keepalive, idles)
}

func HttpClientApp() *FastHttpClient {
	return defaultFastHttpClient
}
