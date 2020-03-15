package consumer

import (
	"fmt"
	"testing"
)

func TestClearConsumer_ConsumeMessage(t *testing.T) {
	consumer, err := ClearConsumer.NewMessageConsumer()
	if nil != err {
		t.Fatal(err)
	}
	err = consumer.ConsumeMessage(msgNil)
	fmt.Println("Err1:", err)
	err = consumer.ConsumeMessage(msgEmpty)
	fmt.Println("Err2:", err)
	err = consumer.ConsumeMessage(msgDefault)
	fmt.Println("Err3:", err)
}
