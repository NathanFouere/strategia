import { WebSocketService } from '@/services/websocket.service'
import { inject, injectable } from 'inversify'
import { usePlayerStore } from '@/stores/player.store'
import type SetInWaitingLobbyPayload from '@/ws-exchange/set-in-waiting-lobby-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import GenerateWsTemplate from '@/utils/generate-ws-template'
import ConnectionHandler from '@/ws-handler/connection.handler'
import WaitingGameHandler from '@/ws-handler/waiting-game.handler'
import RedirectToGameHandler from '@/ws-handler/redirect-to-game.handler'

@injectable()
export default class HomeViewPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    @inject(ConnectionHandler)
    private readonly connectionHandler: ConnectionHandler,
    @inject(WaitingGameHandler)
    private readonly waitingGameHandler: WaitingGameHandler,
    @inject(RedirectToGameHandler)
    private readonly redirectToGameHandler: RedirectToGameHandler,
    public playerStore = usePlayerStore(),
  ) {}

  public initialize(): void {
    this.connectionHandler.subscribe()
    this.waitingGameHandler.subscribe()
    this.redirectToGameHandler.subscribe()
    if (this.playerStore.player) {
      const payload: SetInWaitingLobbyPayload = {
        player_id: this.playerStore.player.id,
      }

      this.wsService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.SET_IN_WAITING_LOBBY, payload))
    }
  }

  public handlePageUnmounted(): void {
    this.connectionHandler.unsubscribe()
    this.waitingGameHandler.unsubscribe()
    this.redirectToGameHandler.unsubscribe()
  }
}
