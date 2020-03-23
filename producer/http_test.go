package producer

import (
	"fmt"
	"github.com/xuzhuoxi/LegoMQ-go/message"
	"github.com/xuzhuoxi/infra-go/eventx"
	"net/http"
	"testing"
	"time"
)

func TestHttpMessageProducer_HttpServer(t *testing.T) {
	s := NewHttpMessageProducer()
	s.MapFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		pageInfo := []byte("The time is: " + time.Now().Format(time.RFC1123))
		pageInfo = append(pageInfo, []byte("/n")...)
		host := r.Host
		msgBody := r.FormValue("msg")
		pageInfo = append(pageInfo, []byte("msg="+msgBody)...)
		w.Write(pageInfo)

		msg := message.NewMessageContext(host, nil, host, host, msgBody)
		s.NotifyMessageProduced(msg)
	})
	s.AddEventListener(EventMessageOnProducer, onHttpProduced)
	err := s.StartHttpListener(":9000")
	if nil != err {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Minute)
}

func onHttpProduced(evt *eventx.EventData) {
	fmt.Println(evt.Data)
}
