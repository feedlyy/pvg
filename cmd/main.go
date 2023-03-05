package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"pvg/config"
	"pvg/controller"
	"pvg/domain"
	"pvg/repository"
	"pvg/service"
	"time"
)

func init() {
	viper.SetConfigFile(`../config/config.json`)
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

	err = db.AutoMigrate(&domain.Users{})
	if err != nil {
		panic(err)
	}

	kafkaUsr := viper.GetString(`kafka.username`)
	kafkaPwd := viper.GetString(`kafka.password`)
	kafkaAddr := viper.GetString(`kafka.address`)
	kafkaRetry := viper.GetInt(`kafka.retry`)
	kafkaTimeout := viper.GetInt(`kafka.timeout`)
	kafkaConfig := config.GetKafkaConfig(kafkaUsr, kafkaPwd, kafkaTimeout, kafkaRetry)
	producers, err := sarama.NewSyncProducer([]string{kafkaAddr}, kafkaConfig)
	if err != nil {
		logrus.Errorf("Unable to create kafka producer got error %v", err)
		return
	}
	defer func() {
		if err := producers.Close(); err != nil {
			logrus.Errorf("Unable to stop kafka producer: %v", err)
			return
		}
	}()

	logrus.Infof("Success create kafka sync-producer")

	kafka := &domain.KafkaProducer{
		Producer: producers,
	}

	timeoutCtx := viper.GetInt(`context.timeout`)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService, time.Duration(timeoutCtx)*time.Second)

	serverAddr := viper.GetString(`server.address`)
	r := gin.Default()
	r.GET("/users", userController.GetUsers)
	r.GET("/user", userController.GetUserByUsername)
	r.POST("/user/register", userController.Create)
	r.PATCH("/user/:id", userController.Update)
	r.DELETE("/user/:id", userController.DeleteUser)
	r.Run(serverAddr)
}
