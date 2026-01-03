import { inject, injectable } from 'inversify'
import { usePlayerStore } from '@/stores/player.store'
import { WebSocketService } from '@/services/websocket.service'
import GenerateWsTemplate from '@/utils/generate-ws-template'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import type UpdatePlayerPseudoPayload from '@/ws-exchange/update-player-pseudo-payload'

@injectable()
export default class PlayerPseudoComponentPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public playerStore = usePlayerStore(),
  ) {}

  public handlePlayerPseudoChange(newPlayerPseudo: string): void {
    if (!this.playerStore.player) return

    const payload: UpdatePlayerPseudoPayload = {
      player_id: this.playerStore.player.id,
      new_pseudo: newPlayerPseudo,
    }

    this.wsService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.UPDATE_PLAYER_PSEUDO, payload))

    this.playerStore.player.pseudo = newPlayerPseudo
  }
}
