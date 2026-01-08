import { WebSocketService } from '@/services/websocket.service'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import type ServerUpdatePayload from '@/ws-exchange/server-update-payload'
import { inject, injectable } from 'inversify'

@injectable()
export default class ServerUpdateHandler {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
  ) {}

  public subscribe(ctx: CanvasRenderingContext2D): void {
    this.wsService.subscribe(
      WS_MESSAGES_TYPE.SERVER_UPDATE_DATAS,
      (payload: ServerUpdatePayload) => {
        for (const update_data of payload.update_datas) {
          ctx.fillStyle = update_data.color
          ctx.fillRect(update_data.x, update_data.y, 1, 1)
        }
      },
    )
  }
}
