import axios from 'axios'
// import express from 'express'
import { customAlphabet } from 'nanoid'

function nanoid() {
  const _nanoid = customAlphabet(
    '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz',
    11
  )
  return _nanoid()
}

async function createProducts() {
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

  const supportLanguage = []
  supportLanguage.push({
    languageCode: 'th',
    id: productTh.id,
  })
  supportLanguage.push({
    languageCode: 'en',
    id: productEn.id,
  })

  const categories = await getCategory()

  const nType = []
  for (const category of categories) {
    const { id, name } = category
    nType.push({
      id,
      name,
    })
  }

  const product = {
    id: nanoid(),
    name: 'แมวเหมียว',
    description: 'แมวเหมียว',
    lifecycleStatus: 'active',
    supportingLanguage: supportLanguage,
    category: nType,
  }

  if (createProductLanguageBody.length) {
    let productLanguageReq = []
    for (const productLanguage of createProductLanguageBody) {
      productLanguageReq.push(createProductLanguage(productLanguage))
    }
    if (productLanguageReq.length) {
      await Promise.all(productLanguageReq)
    }
  }
  const updateCategoryReq = []
  const categoryBody = []
  for (const category of categories) {
    const { id } = category
    const updateCategory = {
      id,
      product: [
        {
          id: product.id,
          name: product.name,
        },
      ],
    }
    updateCategoryReq.push(updateCategoryId(id, updateCategory))
    categoryBody.push(updateCategory)
  }
  if (updateCategoryReq.length) {
    await Promise.all(updateCategoryReq)
  }

  await createProduct(product)
  return {
    productLanguageTh: productTh,
    productLanguageEn: productEn,
    product,
    category: categoryBody,
  }
}

async function updateCategoryId(id, body) {
  try {
    await axios.patch('http://localhost:2566/category/' + id, body)
  } catch (error) {
    console.error(error)
  }
}

async function createProductLanguage(body) {
  try {
    await axios.post('http://localhost:2566/productLanguage', body)
  } catch (error) {
    console.error(error)
  }
}

async function createProduct(body) {
  try {
    await axios.post('http://localhost:2566/products', body)
  } catch (error) {
    console.error(error)
  }
}

async function getCategory() {
  try {
    const { data } = await axios.get(
      'http://localhost:2566/category?lifecycleStatus=active'
    )
    return data
  } catch (error) {
    console.error(error)
  }
}

async function createCategory() {
  try {
    const category = [
      {
        id: nanoid(),
        name: 'cpu',
        version: '1.0',
        lastUpdate: '2023-10-22',
        lifecycleStatus: 'active',
        validFor: {
          startDateTime: '2023-10-22T08:00:00Z',
          endDateTime: '2029-10-23T08:00:00Z',
        },
      },
      {
        id: nanoid(),
        name: 'gpu',
        version: '1.0',
        lastUpdate: '2023-10-22',
        lifecycleStatus: 'active',
        validFor: {
          startDateTime: '2023-10-22T08:00:00Z',
          endDateTime: '2029-10-23T08:00:00Z',
        },
      },
      {
        id: nanoid(),
        name: 'ram',
        version: '1.0',
        lastUpdate: '2023-10-22',
        lifecycleStatus: 'active',
        validFor: {
          startDateTime: '2023-10-22T08:00:00Z',
          endDateTime: '2029-10-23T08:00:00Z',
        },
      },
    ]

    const reqBody = []
    for (const body of category) {
      reqBody.push(axios.post('http://localhost:2566/category', body))
    }
    await Promise.all(reqBody)
    return category
  } catch (error) {
    console.error(error)
  }
}

async function main() {
  const category = await getCategory()
  if (!category) {
    await createCategory()
  }
  const start = performance.now()
  const data = await createProducts()
  console.log(JSON.stringify(data))
  const end = performance.now()
  const duration = end - start
  console.log('createProducts', duration.toFixed(2))
}

main()
