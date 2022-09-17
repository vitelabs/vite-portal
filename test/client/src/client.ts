import axios, { AxiosInstance, AxiosResponse } from "axios"
import WebSocket from "ws"
import { TestConstants } from "./constants"
import { JsonRpcRequest, JsonRpcResponse } from "./types"
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
  error?: WebSocket.ErrorEvent

  constructor(timeout: number, url: string, clientIp?: string, jwtSubject?: string, jwtSecret?: string) {
    super(timeout, clientIp, jwtSubject, jwtSecret)
    this.ws = new WebSocket(url, {
      handshakeTimeout: timeout,
      timeout: timeout,
      headers: this.createHeaders()
    })
    this.ws.onerror = e => {
      this.error = e
    };
  }
}

export class RpcWsClientWrapper {
  client: RpcWsClient
  connected = false
  requests: JsonRpcRequest[] = []
  errors: JsonRpcResponse<any>[] = []

  constructor(client: RpcWsClient) {
    this.client = client
    client.ws.on('open', () => {
      console.log('connected');
      this.connected = true
    });

    client.ws.on('close', (code: number, reason: Buffer) => {
      console.log(`disconnected: ${code} ${reason}`);
      this.connected = false
    });

    client.ws.on('message', (data: Buffer) => {
      const message = JSON.parse(data.toString())
      console.log(message)
      if (message.method) {
        this.requests.push(message)
      } else if (message.error) {
        this.errors.push(message)
      }
    });
  }
}