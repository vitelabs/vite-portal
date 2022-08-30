import { BaseApp } from "./app"
import { CommonUtil } from "./utils"

export class Orchestrator extends BaseApp {
  constructor(url: string) {
    super(url)
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