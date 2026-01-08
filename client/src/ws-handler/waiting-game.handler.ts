import type WaitingGamePayload from '@/ws-exchange/waiting-game-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { WebSocketService } from '@/services/websocket.service'
import { usePendingGameStore } from '@/stores/pending-game.store'
import { inject, injectable } from 'inversify'

@injectable()
export default class WaitingGameHandler {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public pendingGameStore = usePendingGameStore(),
  ) {}

  public subscribe(): void {
    this.wsService.subscribe(WS_MESSAGES_TYPE.WAITING_GAME, (e: WaitingGamePayload) => {
      this.pendingGameStore.setPendingGameId(e.game_id)
      this.pendingGameStore.setSecondsBeforeLaunch(e.seconds_before_launch)
      this.pendingGameStore.setNumberOfWaitingPlayers(e.number_of_waiting_players)
      this.pendingGameStore.setGameLaunching(e.is_game_launching)
      this.pendingGameStore.setSubscribedToGame(e.is_player_waiting_for_game)
    })
  }

  public unsubscribe(): void {
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.WAITING_GAME)
  }
}
