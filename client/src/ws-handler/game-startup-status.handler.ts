import container from '@/container/container'
import { WebSocketService } from '@/services/websocket.service'
import { useOngoingGameStore } from '@/stores/ongoing-game.store'
import type GameStartupStatusPayload from '@/ws-exchange/game-startup-status-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'

export function registerGameStartupStatusHandler() {
  const websocketService = container.get(WebSocketService)
  const ongoingGameStore = useOngoingGameStore()

  websocketService.subscribe(
    WS_MESSAGES_TYPE.GAME_STARTUP_STATUS,
    (e: GameStartupStatusPayload) => {
      ongoingGameStore.setProgressionPercentage(e.progression_percentage)
      ongoingGameStore.setGameStarted(e.game_started)
    },
  )
}
