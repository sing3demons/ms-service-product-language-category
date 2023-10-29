import { Consumer, ConsumerSubscribeTopics, EachBatchPayload, Kafka, EachMessagePayload, logLevel } from 'kafkajs'
import logger from './utils/logger.js'

const groupId = 'category.category,product.product-1'
const clientId = 'product,category.service'

// export function configKafka(broker: string): Kafka {
//   const kafka = new Kafka({
//     clientId,
//     brokers: () => {
//       logger.info('Kafka broker: ' + broker)
//       return [broker!]
//     },
//     requestTimeout: 25000,
//     retry: {
//       factor: 0,
//       multiplier: 4,
//       maxRetryTime: 25000,
//       retries: 10,
//     },
//   })

//   return kafka
// }

export default class KafkaNode {
  public kafka: Kafka
  public constructor(broker: string) {
    this.kafka = new Kafka({
      logLevel: logLevel.INFO,
      clientId,
      brokers: () => {
        logger.info('Kafka broker: ' + broker)
        return [broker!]
      },
      requestTimeout: 25000,
      retry: {
        factor: 0,
        multiplier: 4,
        maxRetryTime: 25000,
        retries: 10,
      },
    })
  }

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
        console.log('function startConsumer')
        logger.error(e.message)
        throw new Error(e.message)
      }
      return consumer
    }
  }
}
