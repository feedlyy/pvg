package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"pvg/config"
	"pvg/domain"
	"pvg/domain/consumer"
	"pvg/helper"
	repoConsumer "pvg/repository/consumer"
	servConsumer "pvg/service/consumer"
)

func init() {
	viper.SetConfigFile(`../../config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbName := viper.GetString(`database.name`)
	dbPort := viper.GetString(`database.port`)
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&domain.Users{}, &consumer.ActivationCodes{})
	if err != nil {
		panic(err)
	}

	kafkaUsr := viper.GetString(`kafka.username`)
	kafkaPwd := viper.GetString(`kafka.password`)
	kafkaAddr := viper.GetString(`kafka.address`)
	kafkaRetry := viper.GetInt(`kafka.retry`)
	kafkaTimeout := viper.GetInt(`kafka.timeout`)
	kafkaConfig := config.GetKafkaConfig(kafkaUsr, kafkaPwd, kafkaTimeout, kafkaRetry)
	consumers, err := sarama.NewConsumer([]string{kafkaAddr}, kafkaConfig)
	if err != nil {
		logrus.Errorf("Error create kakfa consumer got error %v", err)
	}
	defer func() {
		if err := consumers.Close(); err != nil {
			logrus.Fatal(err)
			return
		}
	}()

	kafkaConsumer := &consumer.KafkaConsumer{
		Consumer: consumers,
	}

	// Initialize Mailjet client
	mailJetApi := viper.GetString(`mailjet.api_key`)
	mailJetSecret := viper.GetString(`mailjet.secret_key`)
	mj := mailjet.NewMailjetClient(mailJetApi, mailJetSecret)

	//apiKey := viper.GetString(`sendgrid.apikey`)
	subACRepo := repoConsumer.NewACRepository(db)
	subACService := servConsumer.NewACService(*kafkaConsumer, subACRepo, mj)

	signals := make(chan os.Signal, 1)
	subACService.Process([]string{helper.EmailTopic}, signals)
}
