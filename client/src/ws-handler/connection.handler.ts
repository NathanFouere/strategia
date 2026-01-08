import type ConnectionPayload from '@/ws-exchange/connection-payload'
import { WebSocketService } from '@/services/websocket.service'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { usePlayerStore } from '@/stores/player.store'
import type Player from '@/models/player'
import { inject, injectable } from 'inversify'

@injectable()
export default class ConnectionHandler {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public playerStore = usePlayerStore(),
  ) {}

  public subscribe(): void {
    if (this.playerStore.player) return

    this.wsService.subscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE, (payload: ConnectionPayload) => {
      const player: Player = {
        id: payload.player_id,
        pseudo: payload.player_pseudo,
      }

      this.playerStore.setPlayer(player)
    })
  }

  public unsubscribe(): void {
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE)
  }
}
