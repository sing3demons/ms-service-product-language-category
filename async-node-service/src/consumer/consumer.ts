import { Consumer, Kafka } from 'kafkajs'
import logger from '../utils/logger.js'
import { createCategory, updateCategory } from '../service/category.js'
import { Category } from '../models/category.js'
import TOPIC from '../constant/topic.js'
import { Producer } from '../producer/producer.js'
import { ProductLanguage } from '../models/productLanguage.js'
import { createProductLanguage } from '../service/productLanguage.js'
import { createProduct } from '../service/product.js'
import { Product } from '../models/product.js'

// export default class Consume {
//   public constructor(private consumer: Consumer, private topics: string[]) {}
//   consumeMessage(): void {
//     try {
//       console.log('consumer message')
//       this.consumer.subscribe({ topics: this.topics, fromBeginning: true })
//       this.consumer.run({
//         eachMessage: async ({ topic, message }) => {
//           switch (topic) {
//             case TOPIC.createCategory:
//               if (message.value) {
//                 const req = JSON.parse(message.value.toString()) as Category
//                 const result = await createCategory(req)
//                 if (result === null || result === undefined) {
//                   await producer(TOPIC.createCategoryFailed, result)
//                 } else {
//                   await producer(TOPIC.createCategorySuccess, result)
//                 }
//               }
//               break
//             case TOPIC.updateCategory:
//               if (message.value) {
//                 const req = JSON.parse(message?.value?.toString()) as Category
//                 if (Array.isArray(req.products)) {
//                   const result = await updateCategory(req)
//                   logger.info(JSON.stringify(result))
//                 }
//               }
//               break
//             case TOPIC.createProductLanguage:
//               if (message.value) {
//                 const req = JSON.parse(message?.value?.toString()) as ProductLanguage
//                 const result = await createProductLanguage(req)
//                 if (result === null || result === undefined) {
//                   await producer(TOPIC.createProductLanguageFailed, result)
//                 } else {
//                   await producer(TOPIC.createProductLanguageSuccess, result)
//                 }
//               }
//               break
//             case TOPIC.createProduct:
//               if (message.value) {
//                 const req = JSON.parse(message?.value?.toString()) as Product
//                 const result = await createProduct(req)
//                 if (result === null || result === undefined) {
//                   await producer(TOPIC.createProductFailed, result)
//                 } else {
//                   await producer(TOPIC.createProductSuccess, result)
//                 }
//               }
//               break

//             default:
//               logger.info(`No handler for topic ${topic}`)
//               break
//           }
//         },
//       })
//     } catch (error) {}
//   }
// }

async function consumeMessage(consumer: Consumer) {
  try {
    console.log('consumer message++++>')
    await consumer.run({
      eachMessage: async ({ topic, partition, message }) => {
        console.log({
          topic,
          partition,
          offset: message.offset,
          value: message?.value?.toString(),
        })
        try {
          console.log('consumer : run => ', topic)
          switch (topic) {
            case TOPIC.createCategory:
              if (message.value) {
                const req = JSON.parse(message.value.toString()) as Category
                const result = await createCategory(req)
                // if (result === null || result === undefined) {
                //   await producer(TOPIC.createCategoryFailed, result)
                // } else {
                //   await producer(TOPIC.createCategorySuccess, result)
                // }
              }
              break
            case TOPIC.updateCategory:
              if (message.value) {
                const req = JSON.parse(message?.value?.toString()) as Category
                if (Array.isArray(req.products)) {
                  const result = await updateCategory(req)
                  logger.info(JSON.stringify(result))
                }
              }
              break
            case TOPIC.createProductLanguage:
              if (message.value) {
                const req = JSON.parse(message?.value?.toString()) as ProductLanguage
                const result = await createProductLanguage(req)
                // if (result === null || result === undefined) {
                //   await producer(TOPIC.createProductLanguageFailed, result)
                // } else {
                //   await producer(TOPIC.createProductLanguageSuccess, result)
                // }
              }
              break
            case TOPIC.createProduct:
              if (message.value) {
                const req = JSON.parse(message?.value?.toString()) as Product
                const result = await createProduct(req)
                // if (result === null || result === undefined) {
                //   await producer(TOPIC.createProductFailed, result)
                // } else {
                //   await producer(TOPIC.createProductSuccess, result)
                // }
              }
              break

            default:
              logger.info(`No handler for topic ${topic}`)
              break
          }
        } catch (error) {
          console.log('error consume message')
          console.error(error)
        }
      },
    })
  } catch (error) {
    console.error("Can't consume message", error)
  }
}

export { consumeMessage }
