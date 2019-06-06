package models

import (
	"github.com/Shopify/sarama"
	"github.com/log-collection/sqlserver2kafka/pkg/logging"
	"github.com/log-collection/sqlserver2kafka/pkg/setting"
	"go.uber.org/zap"
)
var (
	producer sarama.AsyncProducer
	err error
)
func KafkaInit() {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V1_1_0_0

	logging.Debug("start make producer")
	//使用配置,新建一个异步生产者
	//producer, err = sarama.NewAsyncProducer([]string{"192.168.56.122:9093"}, config)
	producer, err = sarama.NewAsyncProducer([]string{setting.Config.Kafka.Host}, config)
	if err != nil {
		logging.Error("创建kafka producer失败",zap.Error(err))
		return
	}
	logging.Debug("启动kafka goroutine查看发送消息情况")
	go func(p sarama.AsyncProducer) {
		for {
			select {
			//case suc := <-p.Successes():
			case <-p.Successes():
				//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
				//logging.Debug("kafka发送成功",zap.Int64("offset",suc.Offset),zap.Time("timestamp",suc.Timestamp),zap.Int32("partitions",suc.Partition))
			case fail := <-p.Errors():
				//fmt.Println("err: ", fail.Err)
				logging.Error("kafka发送失败",zap.Error(fail))
			}
		}
	}(producer)
}
func KafkaProducer(sqlData []byte) {

	//循环判断哪个通道发送过来数据.

		// 发送的消息,主题。
		// 注意：这里的msg必须得是新构建的变量，不然你会发现发送过去的消息内容都是一样的，因为批次发送消息的关系。
		msg := &sarama.ProducerMessage{
			Topic: setting.Config.Kafka.Topic,
		}
		//将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(sqlData)
		//fmt.Println(value)
		//使用通道发送
		producer.Input() <- msg
	}
