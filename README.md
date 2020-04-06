# LegoMQ-go

LegoMQ是一个使用go语言编写的消息中间件。具有高并发、高

## 目录

- [背景](#背景)
- [安装&下载](#安装&下载)
- [使用说明](#使用说明)
	- [配置说明](#配置说明)
	- [功能扩展](#功能扩展)
- [示例](#示例)
- [相关仓库](#相关仓库)
- [维护者](#维护者)
- [使用许可](#使用许可)

## 背景

1. 独立的Log服务器、消息处理服务器的需求。
2. 高并发、高吞吐要求。
3. 组件单独可应用。
4. 高可扩展性要求。

## 安装&下载

	`go get -u github.com/xuzhuoxi/LegoMQ-go`

## 使用说明

1. Producer、Queue和Consumer是LegoMQ的三个核心组件。
2. Routing是路由算法组件。
3. Broker是整合了Producer、Queue、Consumer与Routing的一体化组件，是一个完整的消息中间件。
4. Broker的执行模型如下:


### 配置说明

- ProducerSetting<br>
  消息生产者配置，定义如下：<br>
	```
	type ProducerSetting struct {
		Id       string              // 标识
		Mode     ProducerMode        // 消息生产者模式
		LocateId string              // 位置信息
		Http     ProducerSettingHttp // Http消息接收服务器设置
		RPC      ProducerSettingRPC  // RPC消息接收服务器设置
		Sock     netx.SockParams     // Socket消息接收服务器设置
	}
	```
	
	- Id：<br>唯一标识，不能为空字符，多个Producer间Id不能相同。Id不参与路由算法。
	- Mode：<br>Producer的模式，目前已实现HttpProducer、SockProducer、RPCProducer三种。如要扩展自定义模式，请查看[ProducerMode扩展](#ProducerMode扩展)。
	- LocateId：<br>用于标记当前Producer位置信息，多个Producer间LocateId可以相同，LocateId参与路由算法。
	- Http：<br>对应HttpProducer模式下的配置信息。定义如下：<br>
	  
		```
		type ProducerSettingHttp struct {
			Addr    string
		}

		```
		
		- Addr<br>
		  Http服务器监听地址。

	- RPC：<br>对应SockProducer模式下的配置信息。定义如下：<br>
		  
		```
		type ProducerSettingRPC struct {
			Addr    string
		}

		```
		
		- Addr<br>
		  RPC服务器监听地址。
		  
	- Sock：<br>对应RPCProducer模式下的配置信息。
		定义位于infra-go中的[sock.go](https://github.com/xuzhuoxi/infra-go/blob/master/netx/sock.go)中：<br>
		```
		type SockParams struct {
			Network SockNetwork
			// E.g
			// tcp,udp,quic:	127.0.0.1:9999
			LocalAddress string
			// E.g
			// websocket:	ws://127.0.0.1:9999
			// tcp,udp,quic:	127.0.0.1:9999
			RemoteAddress string
		
			// E.g: /,/echo
			WSPattern string
			// E.g: http://127.0.0.1/，最后的"/"必须
			WSOrigin string
			// E.g: ""
			WSProtocol string
		}
		```
		- Network<br>协议类型，请查看[sock.go](https://github.com/xuzhuoxi/infra-go/blob/master/netx/sock.go)。
		- LocalAddress<br>作为服务器的监听地址。
		- 其它请查看[sock.go](https://github.com/xuzhuoxi/infra-go/blob/master/netx/sock.go)。

- ConsumerSetting<br>
	消息消费者配置，定义如下：<br>	
	```
	type ConsumerSetting struct {
		Id      string       // 标识
		Mode    ConsumerMode // 消息生产者模式
		Formats []string     // 格式匹配信息
	
		Log ConsumerSettingLog	//日志记录配置
	}
	```
	
	- Id：<br>唯一标识，不能为空字符，多个Consumer间Id不能相同。Id不参与路由算法。
	- Mode：<br>Consumer的模式，目前已实现ClearConsumer、PrintConsumer、LogConsumer三种。如要扩展自定义模式，请查看[ConsumerMode扩展](#ConsumerMode扩展)。
	- Formats：<br>用于路由的格式匹配信息，具体表达方法因CosumerMode而异。
	- Log：<br>日志配置。定义如下：<br>
		```
		type ConsumerSettingLog struct {
			Level  logx.LogLevel    // 默认日志等级
			Config []logx.LogConfig // 日志配置
		}
		
		```
		
		- Level：<br>默认记录的日志等级。
		- Config：<br>日志配置项目。
		- 详细请查看[infra-go/logx](https://github.com/xuzhuoxi/infra-go/tree/master/logx)。

- QueueSetting
	消息队列配置，定义如下：<br>
	```
	type QueueSetting struct {
		Id   string    // 标识
		Mode QueueMode // 消息队列模式
		Size int       // 队列容量
	
		LocateId string   // 位置信息
		Formats  []string // 格式匹配信息
	}
	```
	
	- Id：<br>唯一标识，不能为空字符，多个Queue间Id不能相同。Id不参与路由算法。
	- Mode：<br>Queue模式，日前已实现ChannelBlockingQueue、ChannelNBlockingQueue、ArrayQueueUnsafe、ArrayQueueSafe四种。如要扩展自定义模式，请查看[QueueMode扩展](#QueueMode扩展)。
	- Size：<br>消息队列允许的最大容量。
	- LocateId：<br>用于标记当前Queue位置信息，多个Queue间LocateId可以相同，LocateId参与路由算法。
	- Formats：<br>用于路由的格式匹配信息，具体表达方法因CosumerMode而异。

- BrokerRoutingSetting
	Broker中Producer、Queue、Consumer间的路由行为模式配置，定义如下：<br>
	
	```
	type BrokerRoutingSetting struct {
		ProducerRouting routing.RoutingMode `json:"PRouting"` // 队列前置路由(生产者 -> 队列)
	
		QueueRouting         routing.RoutingMode `json:"QRouting"`  // 队列后置路由(队列 -> 消费者)
		QueueRoutingDuration time.Duration       `json:"QDuration"` // 队列后置路由时间片
		QueueRoutingQuantity int                 `json:"QQuantity"` // 队列后置处理批量
	}
	```

	- ProducerRouting：<br>由Producer到Queue的路由模式。
	- QueueRouting：<br>由Queue到Consumer的路由模式。
	- QueueRoutingDuration：<br>从Queue提取数据的时间间隔。要求>=0。
	- QueueRoutingQuantity：<br>从Queue提取数据的批量。要求>=1。
	
- BrokerSetting
	Broker配置，定义如下：<br>
	```// Broker配置信息
	type BrokerSetting struct {
		Producers []producer.ProducerSetting // 生产者
		Queues    []queue.QueueSetting       // 队列
		Consumers []consumer.ConsumerSetting // 消息者
		Routing   BrokerRoutingSetting       // 路由
	}
	```

	- Producers：<br>由ProducerSetting组成的数组。运行时用于创建ProducerGroup。
	- Queues：<br>由QueueSetting组成的数组。运行时用于创建QueueGroup。
	- Consumers：<br>由ConsumerSetting组成的数组。运行时用于创建ConsumerGroup。
	- Routing：BrokerRoutingSetting表达的路由行为信息。
	
### 功能扩展

当LegoMQ默认包含的功能不能满足需求时，可以自行扩展功能并加入到LegoMQ中。以下为四4个重要组件的扩展说明。

#### ProducerMode扩展

1. 定义一个新的ProducerMode，建议使用值与现有([查看producer.go](/producer/producer.go))不一致，也可以直接使用CustomizeProducer进行扩展。
2. 创建新模式对应的结构体(假设为"PX")，实现IMessageProducer接口并完成逻辑开发。
3. 如果想同时增加对ProducerSetting的扩展支持，PX应该同时实现IProducerSettingSupport接口。并在InitProducer函数实现完成的初始化行为。
4. 注册：调用`RegisterProducerMode`注册新模式实例的构造行为。

#### ConsumerMode扩展

1. 定义一个新的ConsumerMode，建议使用值与现有([查看consumer.go](/consumer/consumer.go))不一致，也可以直接使用CustomizeConsumer进行扩展。
2. 创建新模式对应的结构体(假设为"CX")，实现IMessageConsumer接口并完成逻辑开发。
3. 如果想同时增加对ConsumerSetting的扩展支持，CX应该同时实现IConsumerSettingSupport接口。并在InitConsumer函数实现完成的初始化行为。
4. 注册：调用`RegisterConsumerMode`注册新模式实例的构造行为。

#### QueueMode扩展

1. 定义一个新的QueueMode，建议使用值与现有([查看queue.go](/queue/queue.go))不一致，也可以直接使用CustomizeQueue进行扩展。
2. 创建新模式对应的结构体(假设为"QX")，实现IMessageContextQueue接口并完成逻辑开发。
3. 注册：调用`RegisterQueueMode`注册新模式实例的构造行为。

#### RoutingMode扩展

1. 定义一个新的RoutingMode，建议使用值与现有([查看routing.go](/routing/routing.go))不一致，也可以直接使用CustomizeRouting进行扩展。
2. 创建新模式对应的结构体(假设为"RX")，实现IRoutingStrategy和IRoutingStrategyConfig接口并完成逻辑开发。
3. 注册：调用`RegisterRoutingStrategy`注册新模式实例的构造行为。

## 示例

- [http-log](examples/httplog)
	
	监听http接收消息并进行log记录的例子。

## 相关仓库

- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)<br>
基础库，整个snail框架中的大部分简单复用的逻辑都抽象到这个基础库中。

- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc)<br>
打包依赖库，主要用于交叉编译

- json-iterator [https://github.com/json-iterator/go](https://github.com/json-iterator/go)<br>
带对应结构体的Json解释库

## 维护者

xuzhuoxi
<xuzhuoxi@gmail.com>
<mailxuzhuoxi@163.com>

## 使用许可

"LegoMQ-go" 基于MIT[License](/LICENSE)开源。


