import { exec } from "child_process"
import { BaseProcess } from "./process"
import { TestContants } from "./constants"
import { Orchestrator } from "./orchestrator"
import { Relayer } from "./relayer"
import { RelayerConfig } from "./types"

export abstract class VitePortal {
  static handleShutdown(p: BaseProcess) {
    process.on("SIGINT", async function () {
      await p.stop()
    })
    process.on("SIGTERM", async function () {
      await p.stop()
    })
    process.on("SIGQUIT", async function () {
      await p.stop()
    })
  }

  public static startCleanup() {
    exec(
      `./start_cleanup.sh`,
      {
        cwd: TestContants.DefaultBinPath
      },
    )
  }

  public static async startOrchestrator(url: string, authUrl: string, timeout: number) {
    const app = new Orchestrator(url, authUrl, timeout)
    this.handleShutdown(app)
    await app.start()
    return app
  }

  public static async startRelayer(config: RelayerConfig, timeout: number) {
    const app = new Relayer(config, timeout)
    this.handleShutdown(app)
    await app.start()
    return app
  }
}