import type ConnectionPayload from '@/ws-exchange/connection-payload'
import { WebSocketService } from '@/services/websocket.service'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { usePlayerStore } from '@/stores/player.store'
import type Player from '@/models/player'
import container from '@/container/container'

export function registerConnectionHandler() {
  const websocketService = container.get(WebSocketService)

  const playerStore = usePlayerStore()

  if (playerStore.player) return

  websocketService.subscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE, (payload: ConnectionPayload) => {
    const player: Player = {
      id: payload.player_id,
      pseudo: payload.player_pseudo,
    }

    playerStore.setPlayer(player)
  })
}
