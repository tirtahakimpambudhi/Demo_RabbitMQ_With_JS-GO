package test

import (
	"fmt"
	"go_test_rabbitmq/config"
	internalConfig "go_test_rabbitmq/internal/config"
	"go_test_rabbitmq/internal/logger"
	"go_test_rabbitmq/message_broker"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	cfg *internalConfig.Config
	log *logrus.Logger
)

func setup(t *testing.T) *amqp091.Connection {
	cfg = config.NewAppConfig()
	log = logger.NewLogger()
	connection, err := message_broker.ConnectionRabbitMQ(cfg.GetString("rabbitmq.cloud"))
	if err != nil {
		failOnError(err, "Failure Connectiion RabbitMQ", t)
	}
	return connection
}

func TestConnection(t *testing.T) {
	conf := config.NewAppConfig()

	testCases := []struct {
		name      string
		url       string
		isError   bool
		typeError error
	}{
		{name: "should be successfully", url: fmt.Sprintf("%s://%s:%s@%s/%s", conf.GetString("rabbitmq.protocol"), conf.GetString("rabbitmq.user"), conf.GetString("rabbitmq.password"), conf.GetString("rabbitmq.host"), conf.GetString("rabbitmq.vhost")), isError: false, typeError: nil},
		// vhost not found
		{name: "should be failure because vhost not found or no access", url: fmt.Sprintf("%s://%s:%s@%s/%s", conf.GetString("rabbitmq.protocol"), conf.GetString("rabbitmq.user"), conf.GetString("rabbitmq.password"), conf.GetString("rabbitmq.host"), "vhost"), isError: true, typeError: amqp091.ErrVhost},
		// username or password wrong
		{name: "should be failure because username or password wrong", url: fmt.Sprintf("%s://%s:%s@%s/%s", conf.GetString("rabbitmq.protocol"), "root", "password", conf.GetString("rabbitmq.host"), conf.GetString("rabbitmq.vhost")), isError: true, typeError: amqp091.ErrCredentials},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%s - %d", testCase.name, i), func(t *testing.T) {
			conn, err := message_broker.ConnectionRabbitMQ(testCase.url)
			require.Equal(t, testCase.isError, err != nil)
			if err != nil {
				require.True(t, testCase.isError)
				if testCase.typeError != nil {
					require.ErrorIs(t, err, testCase.typeError)
				}
				return
			}
			defer func() {
				err := conn.Close()
				require.NoError(t, err)
			}()
		})
	}
}

func failOnError(err error, msg string, t *testing.T) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		t.FailNow()
	}
}

func tearDown(ch *amqp091.Channel, exchange string, queues []string) {
	for _, queue := range queues {
		_, err := ch.QueueDelete(
			queue, // queue name
			false, // ifUnused
			false, // ifEmpty
			false, // noWait
		)
		if err != nil {
			log.Printf("Failed to delete queue %s: %s", queue, err)
		}
	}
	err := ch.ExchangeDelete(
		exchange, // exchange name
		false,    // ifUnused
		false,    // noWait
	)
	if err != nil {
		log.Printf("Failed to delete exchange %s: %s", exchange, err)
	}
}

func consumeMessages(t *testing.T, ch *amqp091.Channel, queue string, expectedBody string, expectedMsgCount int) int {
	msgCount := 0
	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer", t)

	for d := range msgs {
		require.Equal(t, expectedBody, string(d.Body))
		msgCount++
		if msgCount == expectedMsgCount {
			break
		}
	}
	return msgCount
}

func TestRabbitMQExchanges(t *testing.T) {
	conn := setup(t)
	defer func() {
		err := conn.Close()
		require.NoError(t, err)
	}()
	tests := []struct {
		name         string
		exchangeType string
		exchange     string
		queues       []string
		routingKeys  []string
		bindings     []string
		headers      []amqp091.Table
		messages     []string
		expected     []int
	}{
		{
			name:         "DirectExchange",
			exchangeType: amqp091.ExchangeDirect,
			exchange:     "direct_exchange",
			queues:       []string{"direct_queue_1", "direct_queue_2", "direct_queue_3"},
			routingKeys:  []string{"direct_key_1", "direct_key_2", "direct_Key_3"},
			messages:     []string{"Hello Direct Exchange 1", "Hello Direct Exchange 2", "Hello Direct Exchange 2"},
			expected:     []int{1, 1, 1},
		},
		{
			name:         "FanoutExchange",
			exchangeType: amqp091.ExchangeFanout,
			exchange:     "fanout_exchange",
			routingKeys:  []string{""},
			queues:       []string{"fanout_queue_1", "fanout_queue_2"},
			messages:     []string{"Hello Fanout Exchange"},
			expected:     []int{1, 1},
		},
		{
			name:         "TopicExchange",
			exchangeType: amqp091.ExchangeTopic,
			exchange:     "topic_exchange",
			queues:       []string{"topic_queue_1", "topic_queue_2"},
			routingKeys:  []string{"topic.key.1", "topic.key.2"},
			bindings:     []string{"topic.*.1", "topic.*.2"},
			messages:     []string{"Hello Topic Exchange 1", "Hello Topic Exchange 2"},
			expected:     []int{1, 1},
		},
		{
			name:         "HeaderExchange",
			exchangeType: amqp091.ExchangeHeaders,
			exchange:     "header_exchange",
			queues:       []string{"header_queue_1", "header_queue_2"},
			routingKeys:  []string{"", ""},
			headers:      []amqp091.Table{{"key1": "value1"}, {"key2": "value2"}},
			messages:     []string{"Hello Header Exchange 1", "Hello Header Exchange 2"},
			expected:     []int{1, 1},
		},
	}
	for i, tt := range tests {
		i++
		t.Run(fmt.Sprintf("Case %d : %s", i, tt.name), func(t *testing.T) {
			ch, err := conn.Channel()
			failOnError(err, "Failed to open a channel", t)
			defer func() {
				err := ch.Close()
				require.NoError(t, err)
			}()

			err = ch.ExchangeDeclare(
				tt.exchange,     // name
				tt.exchangeType, // type
				true,            // durable
				false,           // auto-deleted
				false,           // internal
				false,           // no-wait
				nil,             // arguments
			)
			// fmt.Println("When im Create Queue and Binding")
			failOnError(err, "Failed to declare an exchange", t)
			for i, queue := range tt.queues {
				queueDeclared, err := ch.QueueDeclare(
					queue, // name
					true,  // durable
					false, // delete when unused
					false, // exclusive
					false, // no-wait
					nil,   // arguments
				)
				t.Logf("Successfully Create Queue %s", queueDeclared.Name)
				failOnError(err, "Failed to declare a queue", t)
				switch tt.exchangeType {
				case amqp091.ExchangeDirect:
					err = ch.QueueBind(
						queue,             //queue name
						tt.routingKeys[i], //routing key
						tt.exchange,       //exchange
						false,
						nil,
					)
				case amqp091.ExchangeHeaders:
					err = ch.QueueBind(
						queue,       // queue name
						"",          // routing key
						tt.exchange, // exchange
						false,
						tt.headers[i],
					)
				case amqp091.ExchangeFanout:
					err = ch.QueueBind(
						queue,       // queue name
						"",          // routing key
						tt.exchange, // exchange
						false,
						nil,
					)
				case amqp091.ExchangeTopic:
					err = ch.QueueBind(
						queue,          // queue name
						tt.bindings[i], // routing key
						tt.exchange,    // exchange
						false,
						nil,
					)
				}
				failOnError(err, "Failed to bind a queue", t)
			}
			// fmt.Println("When im Publish")
			for i, body := range tt.messages {
				publish := amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				}
				// fmt.Println("When i add headers",len(tt.messages),i)
				if tt.exchangeType == amqp091.ExchangeHeaders {
					publish.Headers = tt.headers[i]
				}

				err = ch.Publish(
					tt.exchange,       // exchange
					tt.routingKeys[i], // routing key
					false,             // mandatory
					false,             // immediate
					publish,
				)
				failOnError(err, "Failed to publish a message", t)
			}
			// fmt.Println("When im Consume")
			for i, queue := range tt.queues {
				var actualMessage string
				if len(tt.messages) == 1 {
					actualMessage = tt.messages[0]
				} else {
					actualMessage = tt.messages[i]
				}
				msgCount := consumeMessages(t, ch, queue, actualMessage, tt.expected[i])
				require.Equal(t, tt.expected[i], msgCount)
			}

			tearDown(ch, tt.exchange, tt.queues)
		})
	}
}
