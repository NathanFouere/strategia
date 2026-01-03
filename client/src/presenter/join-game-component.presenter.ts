import { inject, injectable } from 'inversify'
import { usePlayerStore } from '@/stores/player.store'
import { WebSocketService } from '@/services/websocket.service'
import GenerateWsTemplate from '@/utils/generate-ws-template'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { usePendingGameStore } from '@/stores/pending-game.store'
import type GameSubscriptionPayload from '@/ws-exchange/game-subscription-payload'
import type GameUnsubscribePayload from '@/ws-exchange/game-unsubscribe-payload'

@injectable()
export default class JoinGameComponentPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public playerStore = usePlayerStore(),
    public pendingGameStore = usePendingGameStore(),
  ) {}

  public sendSubscriptionToGame(): void {
    if (!this.playerStore.player) {
      throw new Error('Should have a connected player')
    }

    if (this.pendingGameStore.isSubscribedToGame) {
      const payload: GameUnsubscribePayload = {
        player_id: this.playerStore.player.id,
      }

      this.wsService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.GAME_UNSUBSCRIBE, payload))
      return
    }

    const payload: GameSubscriptionPayload = {
      player_id: this.playerStore.player.id,
    }

    this.wsService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.GAME_SUBSCRIPTION, payload))
  }
}
