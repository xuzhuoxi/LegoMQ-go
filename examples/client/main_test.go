// Package main
// Create on 2026/2/25
// @author xuzhuoxi
package client

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/xuzhuoxi/infra-go/netx/httpx"
)

func TestServerMain(t *testing.T) {
	data := make(url.Values)
	data.Add("log", "66666666666")
	httpx.HttpPostForm("http://127.0.0.1:9000/ToLog", data, func(res *http.Response, body *[]byte) {
		fmt.Println(res.StatusCode, string(*body))
	})
	time.Sleep(1 * time.Second)
}
