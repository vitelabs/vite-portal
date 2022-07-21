import axios, { AxiosInstance } from "axios";

export class RpcClient {
  provider: AxiosInstance;

  constructor() {
    this.provider = axios.create({
      timeout: 1000,
    });
  }

  send = async (url: string, method: string, params?: []): Promise<void> => {
    const response = await this.provider.post(url, {
      jsonrpc: "2.0",
      id: 1,
      method,
      params
    });
    return response.data
  }
}