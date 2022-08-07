import axios, { AxiosInstance, AxiosResponse } from "axios"
import { JsonRpcRequest } from "./types"

export class RpcClient {
  requestId: number
  provider: AxiosInstance

  constructor() {
    this.requestId = 0
    this.provider = axios.create({
      timeout: 1000,
    })
  }

  send = async (url: string, method: string, params?: [], id?: number): Promise<AxiosResponse<any, any>> => {
    const response = await this.provider.post(url, this.createJsonRpcRequest(method, params, id), {
      headers: {
        "True-Client-IP": "1.2.3.4"
      }
    })
    return response
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