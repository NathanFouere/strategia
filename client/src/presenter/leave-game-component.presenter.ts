import { WebSocketService } from '@/services/websocket.service'
import { inject, injectable } from 'inversify'
import GenerateWsTemplate from '@/utils/generate-ws-template.ts'
import { useRouter } from 'vue-router'
import { usePlayerStore } from '@/stores/player.store.ts'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template.ts'
import type ExitGamePayload from '@/ws-exchange/exit-game-payload.ts'

@injectable()
export default class LeaveGameComponentPresenter {
  constructor(
    @inject(WebSocketService)
    private readonly wsService: WebSocketService,
    public playerStore = usePlayerStore(),
    public router = useRouter(),
  ) {}

  public leaveGame(gameId: string): void {
    this.wsService.unsubscribe(WS_MESSAGES_TYPE.SERVER_UPDATE_DATAS)
    const exitGamePayload: ExitGamePayload = {
      player_id: this.playerStore.player!.id,
      game_id: gameId,
    }

    this.wsService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.EXIT_GAME, exitGamePayload))

    this.router.push('/')
  }
}
