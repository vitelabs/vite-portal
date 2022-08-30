export type JsonRpcRequest = {
  jsonrpc: string
  id: number
  method: string
  params: any[]
}

export type JsonRpcResponse = {
  jsonrpc: string
  id: number
  result: string
}

export type GenericPage<T> = {
  entries: T[]
  limit: number
  offset: number
  total: number
}

export type NodeEntity = {
  id: string
  chain: string
  rpcHttpUrl: string
  rpcWsUrl: string
}

export type NodeResponse = {
  nodeId: string
  responseTime: number
  response: string
  deadlineExceeded: boolean
  cancelled: boolean
  error: string
}

export type RelayerConfig = {
  rpcUrl: string
  rpcAuthUrl: string
  rpcRelayHttpUrl: string
  repcRelayWsUrl: string
}

export type Relay = {
  host: string
  chain: string
  clientIp: string
  payload: Payload
}

export type RelayResult = {
  sessionKey: string
  relay: Relay
  responses: NodeResponse[]
}

export type Payload = {
  data: string
  method: string
  path: string
  headers: any
}