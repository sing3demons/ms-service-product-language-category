import topic from './constant/topic.js'
import { consumeMessage } from './consumer/consumer.js'
import { connect } from './db.js'
import KafkaNode from './kafka.js'
import dotenv from 'dotenv'
dotenv.config()

async function main() {
  await connect()
  let broker = process.env.KAFKA_BROKERS
  if (!broker) {
    broker = 'kafka:9092'
  }
  console.log(broker)
  const kafkaNode = new KafkaNode(broker)
  const consumer = await kafkaNode.startConsumer()

  const topics: string[] = Object.values(topic).map((item) => item)

  await consumer.subscribe({ topics, fromBeginning: true })
  await consumeMessage(consumer)
}
main().catch((e) => {
  console.log('main error')
  console.log(e)
})
