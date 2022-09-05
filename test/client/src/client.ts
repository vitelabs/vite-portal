import axios, { AxiosInstance, AxiosResponse } from "axios"
import { JsonRpcRequest } from "./types"

export abstract class RpcClient {
  requestId: number
  timeout: number

  constructor(timeout: number) {
    this.requestId = 0
    this.timeout = timeout
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
  provider: AxiosInstance

  constructor(timeout: number) {
    super(timeout)
    this.provider = axios.create({
      timeout: timeout,
    })
  }

  send = async (url: string, method: string, params?: any[], id?: number): Promise<AxiosResponse<any, any>> => {
    const response = await this.provider.post(url, this.createJsonRpcRequest(method, params, id), {
      headers: {
        "True-Client-IP": "1.2.3.4"
      }
    })
    return response
  }
}