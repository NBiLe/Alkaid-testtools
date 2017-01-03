package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type WebController struct {
}

// HttpPostForm http post form
func (ths *WebController) HttpPostForm(address, data string, ret interface{}) error {

	resp, err := http.Post(address,
		"application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return err
	}

	return ths.getResponseDecode(resp, ret)
}

func (ths *WebController) getResponseDecode(resp *http.Response, ret interface{}) error {

	if resp == nil {
		return fmt.Errorf("http.Response is nil")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// fmt.Println("[ HTTPGet().body ]\r\n\t", string(body))
	err = json.Unmarshal(body, ret)
	return err
}
