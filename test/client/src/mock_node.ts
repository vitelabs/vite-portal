import express from 'express'
import * as http from 'http'

export abstract class MockNode {
  app: express.Express
  name?: string
  port?: number
  server?: http.Server

  constructor() {
    this.app = express()
  }

  url = () => {
    return `http://127.0.0.1:${this.port}`
  }

  abstract start(port: number): void

  stop = () => {
    this.server?.close(() => {
      console.log(`[${this.name}] on port ${this.port} has been closed.`)
    })
  }
}

export class DefaultMockNode extends MockNode {
  constructor() {
    super()
    this.name = "DefaultMockNode"
  }

  start = (port: number) => {
    if (this.server) {
      return
    }

    this.port = port

    this.app.get('/', (req, res) => {
      res.send('This is a test response!')
    })

    this.app.post('/', (req, res) => {
      res.send('This is a test response!')
    })

    this.server = this.app.listen(port, () => {
      console.log(`[${this.name}] is listening on port ${port}.`)
    })
  }
}
