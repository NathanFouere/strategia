import { WebSocketService } from '@/services/websocket.service'
import { inject, injectable } from 'inversify'
import { usePlayerStore } from '@/stores/player.store'
import type SetInWaitingLobbyPayload from '@/ws-exchange/set-in-waiting-lobby-payload'
import { registerConnectionHandler } from '@/ws-handler/connection.handler'
import { registerWaitingGameHandler } from '@/ws-handler/waiting-game.handler'
import { registerRedirectToGameHandler } from '@/ws-handler/redirect-to-game.handler'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import GenerateWsTemplate from '@/utils/generate-ws-template'

@injectable()
export default class HomeViewPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public playerStore = usePlayerStore(),
  ) {}

  public handlePageMounted(): void {
    registerConnectionHandler()
    registerWaitingGameHandler()
    registerRedirectToGameHandler()
    if (this.playerStore.player) {
      const payload: SetInWaitingLobbyPayload = {
        player_id: this.playerStore.player.id,
      }

      this.wsService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.SET_IN_WAITING_LOBBY, payload))
    }
  }

  public handlePageUnmounted(): void {
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE)
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.WAITING_GAME)
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME)
  }
}
