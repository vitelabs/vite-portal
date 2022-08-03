import axios, { AxiosInstance, AxiosResponse } from "axios";

export class RpcClient {
  provider: AxiosInstance;

  constructor() {
    this.provider = axios.create({
      timeout: 1000,
    });
  }

  send = async (url: string, method: string, params?: []): Promise<AxiosResponse<any, any>> => {
    const response = await this.provider.post(url, {
      jsonrpc: "2.0",
      id: 1,
      method,
      params
    }, {
      headers: {
        "True-Client-IP": "1.2.3.4"
      }
    });
    return response
  }
}