import { describe, expect, it, beforeAll, afterAll } from "bun:test";
import amqp from "amqplib";

/*eslint logical-assignment-operators: ["error", "always"]*/
/* eslint no-console: ["error", { allow: ["warn", "error","info"] }] */

// Configuration object for RabbitMQ connection
const config = {
  protocol: process.env.MESSAGE_BROKER_PROTOCOL,
  user: process.env.MESSAGE_BROKER_USER,
  password: process.env.MESSAGE_BROKER_PASSWORD,
  host: process.env.MESSAGE_BROKER_HOST,
  port: process.env.MESSAGE_BROKER_PORT,
  vhost: process.env.MESSAGE_BROKER_VIRTUAL_HOST,
  cloud : process.env.MESSAGE_BROKER_CLOUD
};

let connection;


const failOnError = (err, msg) => {
  if (err) {
    console.error(`${msg}: ${err.message}`);
  }
};

// Function to tear down the test setup
const tearDown = async (ch, exchange, queues) => {
  for (const queue of queues) {
    try {
      await ch.deleteQueue(queue);
    } catch (err) {
      console.error(`Failed to delete queue ${queue}: ${err.message}`);
    }
  }
  try {
    await ch.deleteExchange(exchange);
  } catch (err) {
    console.error(`Failed to delete exchange ${exchange}: ${err.message}`);
  }
  await ch.close();
};

// Function to consume messages from a queue and verify them
const consumeMessages = async (ch, queue, expectedBody, expectedMsgCount, done) => {
  let msgCount = 0;
  ch.consume(queue, (msg) => {
    if (msg) {
      const receivedBody = msg.content.toString();
      console.info(`Received message: ${receivedBody}`);
      expect(receivedBody).toBe(expectedBody);
      msgCount++;
      if (msgCount === expectedMsgCount) {
        done(msgCount);
      }
    }
  }, { noAck: true });
};

// Describe the test cases
describe("Testing RabbitMQ Exchanges", () => {
  beforeAll(async () => {
    try {
      connection = await amqp.connect(`${config.cloud}`);
    } catch (err) {
      failOnError(err, "Failed to connect to RabbitMQ");
    }
  });

  afterAll(async () => {
    await connection.close();
  });

  const tests = [
    {
      name: "DirectExchange",
      exchangeType: "direct",
      exchange: "direct_exchange",
      queues: ["direct_queue_1", "direct_queue_2", "direct_queue_3"],
      routingKeys: ["direct_key_1", "direct_key_2", "direct_key_3"],
      messages: ["Hello Direct Exchange 1", "Hello Direct Exchange 2", "Hello Direct Exchange 3"],
      expected: [1, 1, 1]
    },
    {
      name: "FanoutExchange",
      exchangeType: "fanout",
      exchange: "fanout_exchange",
      routingKeys: [""],
      queues: ["fanout_queue_1", "fanout_queue_2"],
      messages: ["Hello Fanout Exchange"],
      expected: [1, 1]
    },
    {
      name: "TopicExchange",
      exchangeType: "topic",
      exchange: "topic_exchange",
      queues: ["topic_queue_1", "topic_queue_2"],
      routingKeys: ["topic.key.1", "topic.key.2"],
      bindings: ["topic.*.1", "topic.*.2"],
      messages: ["Hello Topic Exchange 1", "Hello Topic Exchange 2"],
      expected: [1, 1]
    },
    {
      name: "HeaderExchange",
      exchangeType: "headers",
      exchange: "header_exchange",
      queues: ["header_queue_1", "header_queue_2"],
      routingKeys: ["", ""],
      headers: [{ key1: "value1" }, { key2: "value2" }],
      messages: ["Hello Header Exchange 1", "Hello Header Exchange 2"],
      expected: [1, 1]
    }
  ];

  tests.forEach((tt, i) => {
    it(`Case ${i + 1}: ${tt.name}`, async (done) => {
      const channel = await connection.createChannel();
      try {
        await channel.assertExchange(tt.exchange, tt.exchangeType, { durable: true });
        for (const [idx, queue] of tt.queues.entries()) {
          await channel.assertQueue(queue, { durable: true });
          const bindArgs = tt.exchangeType === "headers" ? tt.headers[idx] : {};
          const routingKey = tt.exchangeType === "headers" ? "" : tt.routingKeys[idx];
          await channel.bindQueue(queue, tt.exchange, routingKey, bindArgs);
        }

        for (const [msgIdx, message] of tt.messages.entries()) {
          const publishOptions = { contentType: "text/plain" };
          if (tt.exchangeType === "headers") {
            publishOptions.headers = tt.headers[msgIdx];
          }
          await channel.publish(tt.exchange, tt.routingKeys[msgIdx] || "", Buffer.from(message), publishOptions);
        }

        for (const [idx, queue] of tt.queues.entries()) {
          await consumeMessages(channel, queue, tt.messages[idx] || tt.messages[0], tt.expected[idx], (msgCount) => {
            expect(msgCount).toBe(tt.expected[idx]);
            if (msgCount === tt.expected[idx]) {
              done();
            }
          });
        }

        await tearDown(channel, tt.exchange, tt.queues);
      } catch (error) {
        failOnError(error, `Error during test case ${tt.name}`);
        await channel.close();
      }
    });
  });
});