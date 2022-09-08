import path from "path"
import { RpcHttpClient } from "./client"
import { BaseProcess } from "./process"
import { CommonUtil, FileUtil, getLocalFileUtil } from "./utils"

export class VuilderNode extends BaseProcess {
  fileUtil: FileUtil
  url?: string
  rpcClient: RpcHttpClient

  constructor(timeout: number) {
    super()
    this.fileUtil = getLocalFileUtil()
    this.rpcClient = new RpcHttpClient(timeout)
  }

  name(): string {
    return "node"
  }

  command(): string {
    return "npx"
  }

  killCommand(): string {
    return ""
  }

  args(): string[] {
    return [
      "vuilder",
      "node",
      "--config",
      "node_config_vuilder.json"
    ]
  }

  initAsync = async (): Promise<void> => {
    let cfg = await this.fileUtil.readFileAsync(path.join(this.binPath, "node_config_vuilder.json"))
    cfg = JSON.parse(cfg)
    const nodeCfg = (cfg.nodes as any)[cfg.defaultNode];
    this.url = nodeCfg.http
    return Promise.resolve()
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    if (CommonUtil.isNullOrWhitespace(this.url)) {
      console.log(`[${this.name()}] error: url is not initialized`)
      process.exit(1)
    }
    const response = await this.rpcClient.send(this.url!, "ledger_getSnapshotChainHeight")
    return response.data?.result > 0
  }
}