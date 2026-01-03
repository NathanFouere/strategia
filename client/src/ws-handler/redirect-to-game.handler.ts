import type RedirectToGamePayload from '@/ws-exchange/redirect-to-game-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { WebSocketService } from '@/services/websocket.service'
import { usePendingGameStore } from '@/stores/pending-game.store'
import router from '@/router'
import container from '@/container/container'

export function registerRedirectToGameHandler() {
  const websocketService = container.get(WebSocketService)
  const pendingGameStore = usePendingGameStore()

  websocketService.subscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME, (e: RedirectToGamePayload) => {
    websocketService.unsubscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE)
    websocketService.unsubscribe(WS_MESSAGES_TYPE.WAITING_GAME)
    websocketService.unsubscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME)

    pendingGameStore.unsetAll()
    router.push(`/game?gameId=${e.game_id}`)
  })
}
