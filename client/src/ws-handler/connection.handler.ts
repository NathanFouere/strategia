import type ConnectionPayload from '@/ws-exchange/connection-payload'
import type { WebSocketService } from '@/services/websocket.service'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { usePlayerStore } from '@/stores/player.store'
import type Player from '@/models/player'

export function registerConnectionHandler(ws: WebSocketService) {
  const playerStore = usePlayerStore()

  if (playerStore.player) return

  ws.subscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE, (payload: ConnectionPayload) => {
    const player: Player = {
      id: payload.player_id,
      pseudo: payload.player_pseudo,
    }

    playerStore.setPlayer(player)
  })
}
