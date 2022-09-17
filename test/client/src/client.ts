import axios, { AxiosInstance, AxiosRequestHeaders, AxiosResponse } from "axios"
import WebSocket from "ws"
import { TestConstants } from "./constants"
import { JsonRpcRequest } from "./types"
import { CommonUtil, JwtUtil } from "./utils"

export abstract class RpcClient {
  requestId: number
  timeout: number
  clientIp?: string
  jwtSubject?: string
  jwtSecret?: string

  constructor(timeout: number, clientIp?: string, jwtSubject?: string, jwtSecret?: string) {
    this.requestId = 0
    this.timeout = timeout
    this.clientIp = clientIp
    this.jwtSubject = jwtSubject
    this.jwtSecret = jwtSecret
  }

  protected createHeaders = (): {
    [x: string]: string;
  } => {
    const headers: {
      [x: string]: string;
    } = {}
    if (!CommonUtil.isNullOrWhitespace(this.clientIp)) {
      headers[TestConstants.DefaultHeaderTrueClientIp] = this.clientIp!
    }
    if (!CommonUtil.isNullOrWhitespace(this.jwtSecret)) {
      const token = JwtUtil.CreateDefaultToken(this.jwtSecret!, this.jwtSubject)
      headers[TestConstants.DefaultHeaderAuthorization] = "Bearer " + token
    }
    return headers
  }

  createJsonRpcRequest = (method: string, params?: any[], id?: number): JsonRpcRequest => {
    if (!id) {
      this.requestId++
      id = this.requestId
    }
    return {
      jsonrpc: "2.0",
      id,
      method,
      params: params ?? []
    }
  }
}

export class RpcHttpClient extends RpcClient {
  http: AxiosInstance

  constructor(timeout: number, clientIp?: string, jwtSubject?: string, jwtSecret?: string) {
    super(timeout, clientIp, jwtSubject, jwtSecret)
    this.http = axios.create({
      timeout: timeout,
    })
  }

  send = async (url: string, method: string, params?: any[], id?: number): Promise<AxiosResponse<any, any>> => {
    const response = await this.http.post(url, this.createJsonRpcRequest(method, params, id), {
      headers: this.createHeaders()
    })
    return response
  }
}

export class RpcWsClient extends RpcClient {
  ws: WebSocket

  constructor(timeout: number, url: string, clientIp?: string, jwtSubject?: string, jwtSecret?: string) {
    super(timeout, clientIp, jwtSubject, jwtSecret)
    this.ws = new WebSocket(url, {
      handshakeTimeout: timeout,
      timeout: timeout,
      headers: this.createHeaders()
    })
  }
}