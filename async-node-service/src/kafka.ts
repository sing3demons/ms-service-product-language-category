import { Consumer, ConsumerSubscribeTopics, EachBatchPayload, Kafka, EachMessagePayload } from 'kafkajs'
import logger from './utils/logger.js'

const groupId = 'category.category,product.product'

export class KafkaConfig {
  static config(): Kafka {
    const kafka = new Kafka({
      clientId: 'category.category.service',
      brokers: ['localhost:9092', '127.0.0.1:9092'],
      retry: {
        initialRetryTime: 1000,
        retries: 3,
      },
    })
    return kafka
  }
}

export function configKafka(): Kafka {
  const kafka = new Kafka({
    clientId: 'category.category.service',
    brokers: ['localhost:9092', '127.0.0.1:9092'],
    retry: {
      initialRetryTime: 1000,
      retries: 3,
    },
  })
  return kafka
}

export default class KafkaNode {
  public constructor(private kafka: Kafka) {}

  async ConnectKafka() {
    const consumer = this.kafka.consumer({ groupId })
    try {
      await consumer.connect()
      logger.info('Kafka connected')
    } catch (e) {
      if (e instanceof Error) {
        logger.error(e.message)
        throw new Error(e.message)
      }
    }
  }

  async startConsumer(): Promise<Consumer> {
    const consumer = this.kafka.consumer({ groupId })
    try {
      if (!consumer) {
        throw new Error('Consumer not created')
      }
      await consumer.connect()

      logger.info('Consumer connected')
      return consumer
    } catch (e) {
      if (e instanceof Error) {
        logger.error(e.message)
        throw new Error(e.message)
      }
      return consumer
    }
  }
}
