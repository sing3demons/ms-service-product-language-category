import { MongoClient } from 'mongodb'
import { Category } from './models/category.js'

const url = process.env.MONGO_URL || 'mongodb://mongodb1:27017,mongodb2:27018,mongodb3:27019/?replicaSet=my-replica-set'
const client = new MongoClient(url)

async function connect() {
  await client.connect()
  client.db('category_microservice_db').collection('category').createIndex({ id: 1 }, { unique: true })
  client.db('product_microservice_db').collection('product').createIndex({ id: 1 }, { unique: true })
  client.db('productLanguage_microservice_db').collection('productLanguage').createIndex({ id: 1 }, { unique: true })
  client.db('productPrice_microservice_db').collection('productPrice').createIndex({ id: 1 }, { unique: true })
  client
    .db('productPriceLanguage_microservice_db')
    .collection('productPriceLanguage')
    .createIndex({ id: 1 }, { unique: true })

  console.log('Connected successfully to server')
}

function disconnect() {
  client.close()
}

function collectionCategory() {
  const db = client.db('category_microservice_db')
  return db.collection<Category>('category')
}

function getCollection<T extends Object>(dbName: string) {
  const collectionName = `${dbName}_microservice_db`
  const db = client.db(collectionName)
  return db.collection<T>(dbName)
}

export { connect, disconnect, getCollection }
