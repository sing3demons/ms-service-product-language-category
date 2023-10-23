import { Request, Response } from 'express'
import Logger from './logger.js'
interface Sensitive {
  [key: string]: any
}

class JSONResponse {
  private static logger = Logger

  static success(req: Request, res: Response, message: string, data?: object) {
    req.body.password && delete req.body.password

    if (data) {
      if (Object.prototype.hasOwnProperty.call(data, 'users')) {
        const users = (data as Sensitive)['users']
        users.forEach((item: Sensitive) => {
          if (Object.prototype.hasOwnProperty.call(data, 'password')) {
            if ((item as Sensitive)['password']) {
              delete (item as Sensitive)['password']
            }
          }
        })
      } else if (typeof data === 'object' && data !== null) {
        if (Object.prototype.hasOwnProperty.call(data, 'password')) {
          if ((data as Sensitive)['password']) {
            delete (data as Sensitive)['password']
          }
        }
      }
    }

    this.logger.info(
      JSON.stringify({
        ip: req.ip,
        method: req.method,
        url: req.url,
        query: req.query,
        body: req.body,
        data: data,
      })
    )

    res.status(200).json({
      code: 200,
      message: message || 'success',
      data: data,
    })
  }

  static create(req: Request, res: Response, message: string, data: object) {
    req.body.password && delete req.body.password
    this.logger.info(
      JSON.stringify({
        ip: req.ip,
        method: req.method,
        url: req.url,
        body: req.body,
        data: data,
      })
    )

    res.status(201).json({
      code: 201,
      message: message || 'created',
      data: data,
    })
  }

  static badRequest(req: Request, res: Response, message: string, data?: object) {
    if (req.body?.password) {
      delete req.body.password
    }

    this.logger.info(
      JSON.stringify({
        ip: req.ip,
        method: req.method,
        url: req.url,
        query: req.query,
        body: req.body,
        data: data,
      })
    )

    res.status(400).json({
      code: 400,
      message: message || 'bad request',
      data: data,
    })
  }

  static notFound(req: Request, res: Response, message: string) {
    this.logger.error(
      JSON.stringify({
        ip: req.ip,
        method: req.method,
        url: req.url,
        query: req.query,
      })
    )

    res.status(404).json({
      code: 404,
      message: message || 'not found',
    })
  }

  static unauthorized(req: Request, res: Response, message: string) {
    this.logger.error(
      JSON.stringify({
        ip: req.ip,
        method: req.method,
        url: req.url,
        query: req.query,
      })
    )

    res.status(401).json({
      code: 401,
      message: message || 'unauthorized',
    })
  }

  static forbidden(req: Request, res: Response, message: string) {
    this.logger.info(req.url, {
      ip: req.ip,
      method: req.method,
      url: req.url,
      query: req.query,
    })
    res.status(403).json({
      code: 403,
      message: message || 'forbidden',
    })
  }

  static serverError(req: Request, res: Response, message?: string, data?: object) {
    this.logger.info(req.url, {
      ip: req.ip,
      method: req.method,
      url: req.url,
      query: req.query,
      data: data,
    })

    res.status(500).json({
      code: 500,
      message: message || 'internal server error',
      data: data,
    })
  }
}

export default JSONResponse
