package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	_L "github.com/fuyaocn/evaluatetools/log"
)

// SocketHttp 收发分开处理的http
type SocketHttp struct {
	conn          net.Conn
	https         bool
	unsecureHttps bool
	address       string
	request       *http.Request
	u             *url.URL
	client        *httputil.ClientConn
	StartSend     int64
	CompleteSend  int64
	Result        string
}

// Init 初始化
func (ths *SocketHttp) Init(address string) (err error) {
	ths.address = address
	ths.u, err = url.Parse(address)
	if err != nil {
		return err
	}

	ths.https = false
	if strings.HasPrefix(address, "https://") {
		ths.https = true
	}

	if ths.https {
		config := &tls.Config{}
		if ths.unsecureHttps {
			config.InsecureSkipVerify = true
		}
		ths.conn, err = tls.Dial("tcp", ths.u.Host, config)
	} else {
		ths.conn, err = net.DialTimeout("tcp", ths.u.Host, 60*time.Second)
	}
	if err != nil {
		return err
	}

	ths.client = httputil.NewClientConn(ths.conn, nil)
	return nil
}

// PostForm post 方式发送
func (ths *SocketHttp) PostForm(data string, header map[string]string) (err error) {
	reader := strings.NewReader(data)
	ths.request, err = http.NewRequest("POST", ths.address, reader)
	if err != nil {
		return err
	}
	for k, v := range header {
		ths.request.Header.Set(k, v)
	}
	ths.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ths.StartSend = time.Now().UnixNano()
	err = ths.client.Write(ths.request)
	return err
}

// Response 读回复
func (ths *SocketHttp) Response(ret interface{}) (err error) {
	resp, err := ths.client.Read(ths.request)
	ths.CompleteSend = time.Now().UnixNano()
	if err != nil {
		return fmt.Errorf(" > Response has error :\r\n%+v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(" > Response ioutil.ReadAll has error :\r\n%+v", err)
	}

	// fmt.Println("[ HTTPGet().body ]\r\n\t", string(body))
	err = json.Unmarshal(body, ret)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("JSON unmarshal response body has error : \r\n %+v\r\nBody String = \r\n[\r\n%s\r\n]\r\n", err, string(body))
	}
	return err
}
