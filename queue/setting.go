package queue

type QueueSetting struct {
	Id   string
	Mode QueueMode
	Size int

	LocateId string
	Formats  []string
}

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
