import { WebSocketService } from '@/services/websocket.service'
import { usePlayerStore } from '@/stores/player.store'
import { inject, injectable } from 'inversify'
import type PixelClickPayload from '@/ws-exchange/pixel-click-payload.ts'
import GenerateWsTemplate from '@/utils/generate-ws-template.ts'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template.ts'

@injectable()
export default class GameComponentPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public playerStore = usePlayerStore(),
  ) {}

  public handleClickOnBoard(gameId: string, x: number, y: number) {
    const pixelClickPayload: PixelClickPayload = {
      id_player: this.playerStore.player!.id,
      game_id: gameId,
      x: x,
      y: y,
    }

    this.wsService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.PIXEL_CLICK_EVT, pixelClickPayload))
  }
}
