package consumer

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"os"
	"pvg/domain/consumer"
	"pvg/helper"
	"time"
)

type acService struct {
	kafka consumer.KafkaConsumer
	repo  consumer.ACRepository
}

func NewACService(k consumer.KafkaConsumer, r consumer.ACRepository) consumer.ACService {
	return &acService{
		kafka: k,
		repo:  r,
	}
}

func (a *acService) Process(topics []string, signals chan os.Signal) {
	var (
		chanMessage   = make(chan *sarama.ConsumerMessage, 256)
		err           error
		partitionList []int32
		id            int
		ac            = consumer.ActivationCodes{}
		expiredAt     = time.Now()
		ctx           = context.Background()
	)

	for _, topic := range topics {
		partitionList, err = a.kafka.Consumer.Partitions(topic)
		if err != nil {
			logrus.Errorf("Sub-Service|Unable to get partition got error %v", err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(a.kafka.Consumer, topic, partition, chanMessage)
		}
	}
	logrus.Infof("Sub-Service|Kafka is consuming....")

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			err = json.Unmarshal(msg.Value, &id)
			if err != nil {
				logrus.Errorf("Sub-Service|Failed marshal, err:%v", err)
				break ConsumerLoop
			}

			// insert into activation code table
			ac.UserID = uint(id)
			ac.Code = uint(helper.GenerateActivationCodes())
			ac.ExpiresAt = expiredAt.Add(time.Hour)
			err = a.repo.Insert(ctx, ac)
			if err != nil {
				break ConsumerLoop
			}

			// send email with activation code

		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		logrus.Errorf("Sub-Service|Unable to consume partition %v got error %v", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			logrus.Errorf("Sub-Service|Unable to close partition %v: %v", partition, err)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}
}
