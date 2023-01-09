import { spawn } from "child_process"
import { TestConstants } from "./constants"
import { BaseProcess } from "./process"

export class Kafka {
  private readonly helper: KafkaHelper
  private readonly zookeper: KafkaZookeper
  private readonly broker: KafkaBroker
  private skipped: boolean

  constructor(timeout: number) {
    this.helper = new KafkaHelper()
    this.zookeper = new KafkaZookeper(timeout)
    this.broker = new KafkaBroker(timeout)
    this.skipped = false
  }

  name(): string {
    return "kafka"
  }

  public async start() {
    const connections = await this.helper.getConnections()
    if (connections > 0) {
      this.skipped = true
      console.error(`[${this.name()}] already started (${connections})`)
      return
    }
    await this.zookeper.start()
    await this.broker.start()
    this.zookeper.process?.stdout.removeAllListeners()
    this.broker.process?.stdout.removeAllListeners()
  }

  public async stop() {
    if (this.skipped) {
      return
    }
    await this.broker.stop()
    await this.broker.kill()
    await this.zookeper.stop()
    await this.zookeper.kill()
  }
}

export class KafkaHelper {
  name(): string {
    return "kafka-helper"
  }

  public async getConnections(): Promise<number> {
    const promise = new Promise<string>((resolve, reject) => {
      const proc = spawn("./kafka/connections.sh",
        [],
        {
          cwd: TestConstants.DefaultBinPath,
        }
      )
      proc.on('error', (error) => {
        console.error(`[${this.name()}] error: ${error}`)
        reject()
      })
      proc.stdout?.on('data', (data) => {
        console.log(`[${this.name()}] stdout: ${data}`)
        resolve(String(data))
      })
      proc.stderr?.on('data', (data) => {
        console.error(`[${this.name()}] stderr: ${data}`)
        reject()
      })
    })
    const result = await promise
    console.log(`[${this.name()}] result: ${result}`)
    if (!result.startsWith("Connections:")) {
      return 0
    }
    try {
      return parseInt(result.replace("Connections: ", "").trim())
    } catch (error) {
      console.log(`[${this.name()}] error: ${error}`)
      return 0
    }
  }
}

export class KafkaZookeper extends BaseProcess {
  private readonly helper: KafkaHelper

  constructor(timeout: number) {
    super(timeout)
    this.helper = new KafkaHelper()
  }

  name(): string {
    return "kafka-zookeeper"
  }

  startCommand(): string {
    return "./kafka/start_zookeeper.sh"
  }

  killCommand(): string {
    return "./kafka/stop_zookeeper.sh"
  }

  startArgs(): string[] {
    return []
  }

  init = async (): Promise<void> => {
    return Promise.resolve()
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    const connections = await this.helper.getConnections()
    console.log(`[${this.name()}] connections: ${connections}`)
    return connections > 0
  }
}

export class KafkaBroker extends BaseProcess {
  private readonly helper: KafkaHelper

  constructor(timeout: number) {
    super(timeout)
    this.helper = new KafkaHelper()
  }

  name(): string {
    return "kafka-broker"
  }

  startCommand(): string {
    return "./kafka/start_broker.sh"
  }

  killCommand(): string {
    return "./kafka/stop_broker.sh"
  }

  startArgs(): string[] {
    return []
  }

  init = async (): Promise<void> => {
    return Promise.resolve()
  }

  isUp = async (): Promise<boolean> => {
    if (this.stopped) {
      process.exit(1)
    }
    const connections = await this.helper.getConnections()
    console.log(`[${this.name()}] connections: ${connections}`)
    return connections > 1
  }
}