import { CreateIndexesOptions, MongoClient } from 'mongodb'
import { Category } from './models/category.js'

const url = process.env.MONGO_URL || 'mongodb://mongodb1:27017,mongodb2:27018,mongodb3:27019/?replicaSet=my-replica-set'
const client = new MongoClient(url)

const options: CreateIndexesOptions = { unique: true, sparse: true, checkKeys: true }

const createIndex = [
  {
    dbName: 'category_microservice_db',
    collectionName: 'category',
    indexSpec: { id: 1 },
    options,
  },
  {
    dbName: 'product_microservice_db',
    collectionName: 'product',
    indexSpec: { id: 1 },
    options,
  },
  {
    dbName: 'productLanguage_microservice_db',
    collectionName: 'productLanguage',
    indexSpec: { id: 1 },
    options,
  },
  {
    dbName: 'productPrice_microservice_db',
    collectionName: 'productPrice',
    indexSpec: { id: 1 },
    options,
  },
  {
    dbName: 'productPriceLanguage_microservice_db',
    collectionName: 'productPriceLanguage',
    indexSpec: { id: 1 },
    options,
  },
]

async function connect() {
  await client.connect()

  for (const { dbName, collectionName, indexSpec, options } of createIndex) {
    const db = client.db(dbName)
    const collection = db.collection(collectionName)
    const index = await collection.createIndex(indexSpec, options)
    console.log(`Index ${index} created on ${dbName}.${collectionName}`)
  }

  console.log('Connected successfully to server')
}

function disconnect() {
  client.close()
}

function getCollection<T extends Object>(dbName: string) {
  const collectionName = `${dbName}_microservice_db`
  const db = client.db(collectionName)
  return db.collection<T>(dbName)
}

export { connect, disconnect, getCollection }
