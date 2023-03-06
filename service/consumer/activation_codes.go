package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/sirupsen/logrus"
	"os"
	"pvg/domain"
	"pvg/domain/consumer"
	"pvg/helper"
	"time"
)

type acService struct {
	kafka   consumer.KafkaConsumer
	repo    consumer.ACRepository
	mailjet *mailjet.Client
}

func NewACService(k consumer.KafkaConsumer, r consumer.ACRepository, c *mailjet.Client) consumer.ACService {
	return &acService{
		kafka:   k,
		repo:    r,
		mailjet: c,
	}
}

func (a *acService) Process(topics []string, signals chan os.Signal) {
	var (
		chanMessage    = make(chan *sarama.ConsumerMessage, 256)
		err            error
		partitionList  []int32
		ac             = consumer.ActivationCodes{}
		expiredAt      time.Time
		ctx            = context.Background()
		activationCode = helper.GenerateActivationCodes()
		usr            = domain.Users{}
		loc            *time.Location
	)
	loc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Errorf("Sub-Service|Err when get location %v", err)
	}

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
			err = json.Unmarshal(msg.Value, &usr)
			if err != nil {
				logrus.Errorf("Sub-Service|Failed marshal, err:%v", err)
				break ConsumerLoop
			}

			// insert into activation code table
			ac.UserID = usr.ID
			ac.Code = uint(activationCode)
			expiredAt = time.Now().In(loc)
			ac.ExpiresAt = expiredAt.Add(time.Hour)
			err = a.repo.Insert(ctx, ac)
			if err != nil {
				break ConsumerLoop
			}

			// send email with activation code
			err = a.sendEmail(usr, activationCode)
			if err != nil {
				break ConsumerLoop
			}
			logrus.Println("Email Sent!")
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

func (a *acService) sendEmail(user domain.Users, codes int) error {
	var (
		err     error
		message *mailjet.InfoSendMail
	)

	// Prepare the message
	message = &mailjet.InfoSendMail{
		FromEmail: "octopus.mailtest@gmail.com",
		FromName:  "Octopus Mail",
		Recipients: []mailjet.Recipient{
			{
				Email: user.Email,
			},
		},
		Subject:  "Registration Step",
		TextPart: "",
		HTMLPart: fmt.Sprintf("Hello, this is verification code for completing registration process: <b>%v</b>",
			codes),
	}

	_, err = a.mailjet.SendMail(message)
	if err != nil {
		logrus.Errorf("Sub-Service|Failed Send Email, err:%v", err)
		return err
	}

	return nil
}
