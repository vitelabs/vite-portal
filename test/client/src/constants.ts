import path from "path"

export abstract class TestConstants {
  static DefaultBinPath = path.join(path.dirname(__dirname), "bin")
  static HeaderTrueClientIp = "CF-Connecting-IP"
}