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
    const response = await this.rpcClient.send(this.url, "core_getAppInfo")
    return response.data?.result?.name === "vite-portal-orchestrator"
  }
}