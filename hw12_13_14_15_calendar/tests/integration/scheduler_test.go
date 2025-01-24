package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/model"
	queue "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
)

type QueueMessage interface {
	Send(msg string) error
	Receive() (<-chan string, error)
}

type SchedulerSuite struct {
	suite.Suite
	q  QueueMessage
	ch *amqp.Channel
}

func TestServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(SchedulerSuite))
}

func (s *SchedulerSuite) SetupSuite() {
	cfg := config.Config{
		RabbitMQ: config.RabbitMQ{
			Host:     "rabbitmq",
			Port:     "5672",
			Username: "guest",
			Password: "guest",
		},
	}
	eventQueue, err := queue.NewQueue(&cfg)
	s.Require().NoError(err)

	s.q = eventQueue
	s.ch = eventQueue.Channel
}

func (s *SchedulerSuite) TestSendMessage() {
	n := model.Notification{
		EventID: uuid.MustParse("53aa35c8-e659-44b2-882f-f6056e443c99"),
		Title:   "notification title",
		Date:    time.Now(),
		UserID:  "1000",
	}
	msg := fmt.Sprintf("Notification to User: %s, Event ID: %s, Title: %s, Notify At: %s",
		n.UserID, n.EventID, n.Title, n.Date)

	err := s.q.Send(msg)
	s.Require().NoError(err)
	messages, err := s.ch.Consume(
		"notifications",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	s.Require().NoError(err)
	for msg := range messages {
		fmt.Printf("Received message: %s\n", string(msg.Body))
		s.Require().NotEmpty(msg)
		break
	}
}
