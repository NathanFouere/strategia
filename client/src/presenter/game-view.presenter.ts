import { WebSocketService } from '@/services/websocket.service'
import { useOngoingGameStore } from '@/stores/ongoing-game.store'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { registerGameStartupStatusHandler } from '@/ws-handler/game-startup-status.handler'
import { inject, injectable } from 'inversify'

@injectable()
export default class GameViewPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public ongoingGameStore = useOngoingGameStore(),
  ) {}

  public handlePageMounted(): void {
    registerGameStartupStatusHandler()
  }

  public handlePageUnmounted(): void {
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.GAME_STARTUP_STATUS)
  }
}
