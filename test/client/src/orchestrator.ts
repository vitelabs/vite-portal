import { BaseApp } from "./app"
import { CommonUtil } from "./utils"
import { GenericPage, RelayerEntity } from "../src/types"

export class Orchestrator extends BaseApp {
  authUrl: string

  constructor(url: string, authUrl: string, timeout: number) {
    super(url, timeout)
    this.authUrl = authUrl
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

  getRelayers = async (offset?: number, limit?: number): Promise<GenericPage<RelayerEntity>> => {
    const params = [
      !!offset ? offset : 0,
      !!limit ? limit : 0
    ]
    const response = await this.rpcClient.send(this.authUrl, "relayers_getPaginated", params)
    return response.data.result
  }
}