package queue

// QueueSetting
// 消息队列设置
type QueueSetting struct {
	Id   string    // 标识
	Mode QueueMode // 消息队列模式
	Size int       // 队列容量

	LocateId string   // 位置信息
	Formats  []string // 格式匹配信息
}

// NewMessageQueue
// 创建消息队列
// err:
//
//	ErrQueueModeUnregister: 实例化功能未注册
func (qs QueueSetting) NewMessageQueue() (q IMessageContextQueue, err error) {
	q, err = qs.Mode.NewContextQueue(qs.Size)
	if nil != err {
		return nil, err
	}
	q.SetId(qs.Id)
	q.SetLocateId(qs.LocateId)
	q.SetFormats(qs.Formats)
	return
}
