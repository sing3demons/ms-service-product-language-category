import { Collection, MongoClient } from 'mongodb'
import { Category } from './models/category.js'
import { ProductLanguage } from './models/productLanguage.js'
const url = 'mongodb://root:category1234@127.0.0.1:27017/?authSource=admin'
const client = new MongoClient(url)

const dbName = 'category_microservice_db'
const collectionName = 'category'

async function connect() {
  await client.connect()
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
