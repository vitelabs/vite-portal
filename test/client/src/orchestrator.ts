import { BaseApp } from "./app";
import { RpcClient } from "./client";
import { CommonUtil } from "./utils";

export class Orchestrator extends BaseApp {
  rpcClient: RpcClient

  constructor(url: string) {
    super(url)
    this.rpcClient = new RpcClient()
  }

  name(): string {
    return "orchestrator"
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    const request = this.rpcClient.createJsonRpcRequest("public_version")
    const response = await this.axiosClient.post("/", request)
    return !CommonUtil.isNullOrWhitespace(response.data.result)
  }
}