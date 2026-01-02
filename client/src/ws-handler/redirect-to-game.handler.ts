import type RedirectToGamePayload from '@/ws-exchange/redirect-to-game-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import type { WebSocketService } from '@/services/websocket.service'
import { usePendingGameStore } from '@/stores/pending-game.store'
import router from '@/router'

export function registerRedirectToGameHandler(ws: WebSocketService) {
  const pendingGameStore = usePendingGameStore()

  ws.subscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME, (e: RedirectToGamePayload) => {
    ws.unsubscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE)
    ws.unsubscribe(WS_MESSAGES_TYPE.WAITING_GAME)
    ws.unsubscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME)

    pendingGameStore.unsetAll()
    router.push(`/game?gameId=${e.game_id}`)
  })
}
