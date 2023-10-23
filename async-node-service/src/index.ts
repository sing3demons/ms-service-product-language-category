import { Kafka } from 'kafkajs'
import topic from './constant/topic.js'
import { consumeMessage } from './consumer/consumer.js'
import { connect, disconnect } from './db.js'
import KafkaNode from './kafka.js'

const kafka = new Kafka({
  clientId: 'category.category.service',
  brokers: ['localhost:9092', '127.0.0.1:9092'],
  retry: {
    initialRetryTime: 1000,
    retries: 3,
  },
})

async function main() {
  await connect()
  const kafkaNode = new KafkaNode(kafka)
  const consumer = await kafkaNode.startConsumer()

  const topics: string[] = Object.values(topic).map((item) => item)

  await consumeMessage(consumer, topics)
}
main().catch(console.error)
