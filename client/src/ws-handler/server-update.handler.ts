import { WebSocketService } from '@/services/websocket.service'
import container from '@/container/container'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import type ServerUpdatePayload from '@/ws-exchange/server-update-payload'

export function registerServerUpdateHandler(ctx: CanvasRenderingContext2D) {
  const websocketService = container.get(WebSocketService)

  websocketService.subscribe(
    WS_MESSAGES_TYPE.SERVER_UPDATE_DATAS,
    (payload: ServerUpdatePayload) => {
      for (const update_data of payload.update_datas) {
        ctx.fillStyle = update_data.color
        ctx.fillRect(update_data.x, update_data.y, 1, 1)
      }
    },
  )
}
