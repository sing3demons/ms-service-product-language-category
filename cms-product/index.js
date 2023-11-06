import axios from 'axios'
// import express from 'express'
import { fakerDE as faker } from '@faker-js/faker'
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
  const productTh = createProductsLanguage('th')
  createProductLanguageBody.push(productTh)
  const productEn = createProductsLanguage('en')
  createProductLanguageBody.push(productEn)
  const productMy = createProductsLanguage('my')
  createProductLanguageBody.push(productMy)
  const productKm = createProductsLanguage('km')
  createProductLanguageBody.push(productKm)

  const supportLanguage = []
  supportLanguage.push({
    languageCode: 'th',
    id: productTh.id,
    '@referredType': 'products',
  })
  supportLanguage.push({
    languageCode: 'en',
    id: productEn.id,
    '@referredType': 'products',
  })
  supportLanguage.push({
    languageCode: 'my',
    id: productMy.id,
    '@referredType': 'products',
  })
  supportLanguage.push({
    languageCode: 'km',
    id: productKm.id,
    '@referredType': 'products',
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
    name: faker.commerce.productName(),
    description: faker.commerce.productDescription(),
    version: '1.0',
    lastUpdate: new Date().toISOString(),
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
      products: [
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

function createProductsLanguage(langCode) {
  const body = {
    id: nanoid(),
    languageCode: langCode.toLowerCase(),
    attachment: [
      {
        id: nanoid(),
        attachmentType: 'image',
        mimeType: 'jpg',
        name: faker.commerce.productName(),
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
  return body
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
  const categories = await getCategory()
  if (Array.isArray(categories) && categories.length === 0) {
    await createCategory()
    return
  }
  const start = performance.now()
  const { productLanguageTh, productLanguageEn, product, category } =
    await createProducts()
  console.log(JSON.stringify(productLanguageTh))
  console.log(JSON.stringify(productLanguageEn))
  console.log(JSON.stringify(product))
  console.log(JSON.stringify(category))
  const end = performance.now()
  const duration = end - start
  console.log('createProducts', duration.toFixed(2))
}

main()
