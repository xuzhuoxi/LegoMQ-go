package message

type IMessageContext interface {
	Header() interface{}

	Sender() string
	Receiver() string

	Timestamp() int64
	Index() uint64

	Body() interface{}
}
