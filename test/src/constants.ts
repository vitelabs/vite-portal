import path from "path"

export abstract class TestConstants {
  static DefaultBinPath = path.join(path.dirname(__dirname), "bin")
  static DefaultChain = "chain1"
  static DefaultJwtSecret = "secret1234"
  static DefaultJwtRelayerIssuer = "vite-portal-relayer"
  static DefaultHeaderAuthorization = "Authorization"
  static DefaultHeaderTrueClientIp = "CF-Connecting-IP"
  static DefaultIpAddress = "0.0.0.0"
  static DefaultRpcNodeTimeout = 2000
  static DefaultPageLimit = 1000
  static SupportedChains = {
    ViteMain: "vite_main",
    ViteBuidl: "vite_buidl"
  }
}