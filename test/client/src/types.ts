export type JsonRpcResponse = {
  jsonrpc?: string
  id?: number
  result?: string
}

export type NodeEntity = {
  id: string
  chain: string
  rpcHttpUrl: string
  rpcWsUrl: string
}

export type GenericPage<T> = {
  entries: T[]
  limit: number
  offset: number
  total: number
}