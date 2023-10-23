import { Consumer, Kafka } from 'kafkajs'
import logger from '../utils/logger.js'
import { createCategory, updateCategory } from '../service/category.js'
import { Category } from '../models/category.js'
import TOPIC from '../constant/topic.js'
import { producer } from '../producer/producer.js'
import { ProductLanguage } from '../models/productLanguage.js'
import { createProductLanguage } from '../service/productLanguage.js'
import { createProduct } from '../service/product.js'
import { Product } from '../models/product.js'

const groupId = 'category.category'
const kafka = new Kafka({
  clientId: 'category-app',
  brokers: ['localhost:9092', '127.0.0.1:9092'],
})
async function startConsumer(): Promise<Consumer> {
  const consumer = kafka.consumer({ groupId })
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

async function consumeMessage(consumer: Consumer, topics: string[]) {
  await consumer.subscribe({ topics, fromBeginning: true })
  await consumer.run({
    eachMessage: async ({ topic, message }) => {
      if (topic === TOPIC.createCategory && message.value) {
        const req = JSON.parse(message?.value?.toString()) as Category
        const result = await createCategory(req)
        if (result === null || result === undefined) {
          await producer(TOPIC.createCategoryFailed, result)
        } else {
          await producer(TOPIC.createCategorySuccess, result)
        }
      } else if (topic === TOPIC.updateCategory && message.value) {
        const req = JSON.parse(message?.value?.toString()) as Category
        if (Array.isArray(req.products)) {
          const result = await updateCategory(req)
          logger.info(JSON.stringify(result))
        }
      } else if (topic === TOPIC.createProductLanguage && message.value) {
        const req = JSON.parse(message?.value?.toString()) as ProductLanguage
        const result = await createProductLanguage(req)
        if (result === null || result === undefined) {
          await producer(TOPIC.createProductLanguageFailed, result)
        } else {
          await producer(TOPIC.createProductLanguageSuccess, result)
        }
      } else if (topic === TOPIC.createProduct && message.value) {
        const req = JSON.parse(message?.value?.toString()) as Product
        const result = await createProduct(req)
        if (result === null || result === undefined) {
          await producer(TOPIC.createProductFailed, result)
        } else {
          await producer(TOPIC.createProductSuccess, result)
        }
      }
    },
  })
}

export { startConsumer, consumeMessage }
