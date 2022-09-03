import path from "path"

export abstract class TestContants {
  static DefaultBinPath = path.join(path.dirname(__dirname), "bin")
}