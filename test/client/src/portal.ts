import { exec } from "child_process"
import { BaseApp } from "./app"
import { TestContants } from "./constants"
import { Orchestrator } from "./orchestrator"
import { Relayer } from "./relayer"
import { RelayerConfig } from "./types"

export abstract class VitePortal {
  static handleShutdown(app: BaseApp) {
    process.on("SIGINT", async function () {
      await app.stop()
    })
    process.on("SIGTERM", async function () {
      await app.stop()
    })
    process.on("SIGQUIT", async function () {
      await app.stop()
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