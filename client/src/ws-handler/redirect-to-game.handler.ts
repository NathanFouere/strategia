import type RedirectToGamePayload from '@/ws-exchange/redirect-to-game-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { WebSocketService } from '@/services/websocket.service'
import { usePendingGameStore } from '@/stores/pending-game.store'
import router from '@/router'
import { inject, injectable } from 'inversify'

@injectable()
export default class RedirectToGameHandler {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public pendingGameStore = usePendingGameStore(),
  ) {}

  public subscribe(): void {
    this.wsService.subscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME, (e: RedirectToGamePayload) => {
      this.wsService.unsubscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE)
      this.wsService.unsubscribe(WS_MESSAGES_TYPE.WAITING_GAME)
      this.wsService.unsubscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME)

      this.pendingGameStore.unsetAll()
      router.push(`/game?gameId=${e.game_id}`)
    })
  }

  public unsubscribe(): void {
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME)
  }
}
