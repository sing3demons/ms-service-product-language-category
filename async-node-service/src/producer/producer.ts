import { Kafka } from 'kafkajs'
import logger from '../utils/logger.js'
const kafka = new Kafka({
  clientId: 'category.category.service',
  brokers: ['localhost:9092', '127.0.0.1:9092'],
})

async function producer(topic: string, value: any) {
  const producer = kafka.producer()
  await producer.connect()
  logger.info(
    JSON.stringify({
      topic,
      value,
    })
  )

  await producer.send({
    topic,
    messages: [{ value: JSON.stringify(value) }],
  })

  await producer.disconnect()
}

export { producer }
