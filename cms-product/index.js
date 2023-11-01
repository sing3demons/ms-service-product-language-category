import axios from 'axios'
import express from 'express'
import { customAlphabet } from 'nanoid'

function nanoid() {
   const _nanoid =  customAlphabet(
    '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz',
    11
  )
  return _nanoid()
}

async function createProduct() {
    const createProductLanguageBody = []
  const productTh = {
    id: nanoid(),
    languageCode: 'th',
    attachment: [
      {
        id: nanoid(),
        attachmentType: 'image',
        mimeType: 'jpg',
        url: 'https://randomuser.me/api/portraits/men/11',
        redirectUrl: 'https://www.google.com',
        validFor: {
          endDateTime: '2023-10-25T23:59:59+07:00',
          startDateTime: '2023-09-28T00:00:00+07:00',
        },
        displayInfo: {
          valueType: 'plaintext',
          value: ['แมวเหมียว'],
        },
      },
    ],
  }
  createProductLanguageBody.push(productTh)

  const productEn = {
    id: nanoid(),
    languageCode: 'en',
    attachment: [
      {
        id: nanoid(),
        attachmentType: 'image',
        mimeType: 'jpg',
        url: 'https://randomuser.me/api/portraits/men/12',
        redirectUrl: 'https://www.google.com',
        validFor: {
          endDateTime: '2023-10-25T23:59:59+07:00',
          startDateTime: '2023-09-28T00:00:00+07:00',
        },
        displayInfo: {
          valueType: 'plaintext',
          value: ['cat meow'],
        },
      },
    ],
  }
    createProductLanguageBody.push(productEn)

    const category = [
      {
        id: nanoid(),
        name: 'cpu',
        version: '1.0.0',
        lastUpdate: '2023-10-22',
        validFor: {
          startDateTime: '2023-10-22T08:00:00Z',
          endDateTime: '2029-10-23T08:00:00Z',
        },
      },
      {
        id: nanoid(),
        name: 'gpu',
        version: '1.0.0',
        lastUpdate: '2023-10-22',
        validFor: {
          startDateTime: '2023-10-22T08:00:00Z',
          endDateTime: '2029-10-23T08:00:00Z',
        },
      },
      {
        id: nanoid(),
        name: 'ram',
        version: '1.0.0',
        lastUpdate: '2023-10-22',
        validFor: {
          startDateTime: '2023-10-22T08:00:00Z',
          endDateTime: '2029-10-23T08:00:00Z',
        },
      },
    ]

  const product = {}

  if (Array.isArray(createProductLanguageBody)&& createProductLanguageBody.length > 0) {
    // product.productLanguage = createProductLanguageBody
    let productLanguageReq = []
    for (const productLanguage of createProductLanguageBody) {
      productLanguageReq.push(createProductLanguage(productLanguage))
    }
    if (productLanguageReq.length){
        await Promise.all(productLanguageReq)
    }
  }
 
}

async function createProductLanguage(data) {
  try {
    const { data } = axios.post('http://localhost:2566/productLanguage', data)
  } catch (error) {
    console.error(error)
  }
}

async function main() {
    const id = nanoid()
    console.log(id)
//   const app = express()
//   createProduct()
//   app.listen(3000, () => console.log('Server running on port 3000'))
}

main()
