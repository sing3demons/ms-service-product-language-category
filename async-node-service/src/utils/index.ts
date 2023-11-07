import { KafkaMessage } from 'kafkajs'
import logger from './logger.js'

export function GetDataFromEvent<T>(message: KafkaMessage): T {
  const event = {
    offset: message.offset,
    data: message?.value?.toString(),
    key: message?.key?.toString(),
    headers: message.headers,
    timestamp: message.timestamp,
    attributes: message.attributes,
    size: message.size,
  }

  const headers = event.headers
  const body: Event = event?.data ? JSON.parse(event.data) : null
  const data = body.body

  logger.info('event ==================>')
  logger.info(JSON.stringify(event))

  console.log('header', JSON.stringify(headers))
  if (typeof data == 'string') {
    return JSON.parse(data) as T
  } else {
    return data as T
  }
}

interface Event {
  header: Header
  body: any
}

interface Header {
  version?: string
  timestamp?: string
  orgService?: string
  from?: string
  channel?: string
  broker?: string
  session?: string
  transaction?: string
  communication?: string
  groupTags?: any[]
  identity?: Identity
  baseApiVersion?: string
  schemaVersion?: string
  instanceData?: string
}

interface Identity {
  device: number
}
