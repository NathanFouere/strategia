export interface ServerUpdate {
  x: number
  y: number
  color: string
}

export default interface ServerUpdatePayload {
  update_datas: ServerUpdate[]
}
