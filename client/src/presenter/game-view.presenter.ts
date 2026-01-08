import { WebSocketService } from '@/services/websocket.service'
import { useOngoingGameStore } from '@/stores/ongoing-game.store'
import GameStartupStatusHandler from '@/ws-handler/game-startup-status.handler'
import { inject, injectable } from 'inversify'

@injectable()
export default class GameViewPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    @inject(GameStartupStatusHandler)
    private readonly gameStartupStatusHandler: GameStartupStatusHandler,
    public ongoingGameStore = useOngoingGameStore(),
  ) {}

  public initialize(): void {
    this.gameStartupStatusHandler.subscribe()
  }

  public handlePageUnmounted(): void {
    this.gameStartupStatusHandler.unsubscribe()
  }
}
