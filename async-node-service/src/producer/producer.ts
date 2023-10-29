import { Kafka } from 'kafkajs'
import logger from '../utils/logger.js'

let broker = process.env.KAFKA_BROKERS
if (!broker) {
  broker = 'kafka:9092'
}

const kafka = new Kafka({
  clientId: 'product,category.service',
  brokers: () => {
    logger.info('Kafka broker: ' + broker)
    return [broker!]
  },
})

async function Producer(topic: string, value: any) {
  try {
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
  } catch (error) {
    console.log('producer error')
    console.error(error)
  }
}

export { Producer }
