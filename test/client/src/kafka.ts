import { BaseProcess } from "./process"

export class KafkaZookeper extends BaseProcess {
  constructor(timeout: number) {
    super(timeout)
  }

  name(): string {
    return "kafka-zookeeper"
  }

  startCommand(): string {
    return "./start_zookeeper.sh"
  }

  killCommand(): string {
    return "./stop_zookeeper.sh"
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
    return true
  }
}