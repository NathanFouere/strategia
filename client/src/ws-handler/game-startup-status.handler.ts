import { WebSocketService } from '@/services/websocket.service'
import { useOngoingGameStore } from '@/stores/ongoing-game.store'
import type GameStartupStatusPayload from '@/ws-exchange/game-startup-status-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { inject, injectable } from 'inversify'

@injectable()
export default class GameStartupStatusHandler {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public ongoingGameStore = useOngoingGameStore(),
  ) {}

  public subscribe(): void {
    this.wsService.subscribe(
      WS_MESSAGES_TYPE.GAME_STARTUP_STATUS,
      (e: GameStartupStatusPayload) => {
        this.ongoingGameStore.setProgressionPercentage(e.progression_percentage)
        this.ongoingGameStore.setGameStarted(e.game_started)
      },
    )
  }
  public unsubscribe(): void {
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.GAME_STARTUP_STATUS)
  }
}
