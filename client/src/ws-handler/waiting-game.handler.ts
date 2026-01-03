import type WaitingGamePayload from '@/ws-exchange/waiting-game-payload'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import { WebSocketService } from '@/services/websocket.service'
import { usePendingGameStore } from '@/stores/pending-game.store'
import container from '@/container/container'

export function registerWaitingGameHandler() {
  const websocketService = container.get(WebSocketService)
  const pendingGameStore = usePendingGameStore()

  websocketService.subscribe(WS_MESSAGES_TYPE.WAITING_GAME, (e: WaitingGamePayload) => {
    pendingGameStore.setPendingGameId(e.game_id)
    pendingGameStore.setSecondsBeforeLaunch(e.seconds_before_launch)
    pendingGameStore.setNumberOfWaitingPlayers(e.number_of_waiting_players)
    pendingGameStore.setGameLaunching(e.is_game_launching)
    pendingGameStore.setSubscribedToGame(e.is_player_waiting_for_game)
  })
}
