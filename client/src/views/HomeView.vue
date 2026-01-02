<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue'
import { usePlayerStore } from '@/stores/player.store'
import { usePendingGameStore } from '@/stores/pending-game.store'
import container from '@/container/container'
import { WebSocketService } from '@/services/websocket.service'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template'
import GenerateWsTemplate from '@/utils/generate-ws-template'
import type GameSubscriptionPayload from '@/ws-exchange/game-subscription-payload'
import type GameUnsubscribePayload from '@/ws-exchange/game-unsubscribe-payload'
import type SetInWaitingLobbyPayload from '@/ws-exchange/set-in-waiting-lobby-payload'
import type UpdatePlayerPseudoPayload from '@/ws-exchange/update-player-pseudo-payload'
import { registerConnectionHandler } from '@/ws-handler/connection.handler'
import { registerWaitingGameHandler } from '@/ws-handler/waiting-game.handler'
import { registerRedirectToGameHandler } from '@/ws-handler/redirect-to-game.handler'
const playerStore = usePlayerStore()
const pendingGameStore = usePendingGameStore()
const websocketService = container.get(WebSocketService)

registerConnectionHandler(websocketService)
registerWaitingGameHandler(websocketService)
registerRedirectToGameHandler(websocketService)

if (playerStore.player) {
  const payload: SetInWaitingLobbyPayload = {
    player_id: playerStore.player.id,
  }

  websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.SET_IN_WAITING_LOBBY, payload))
}

function sendSubscriptionToGame(): void {
  if (!playerStore.player) {
    throw new Error('Should have a connected player')
  }

  if (pendingGameStore.isSubscribedToGame) {
    const payload: GameUnsubscribePayload = {
      player_id: playerStore.player.id,
    }

    websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.GAME_UNSUBSCRIBE, payload))
    return
  }

  const payload: GameSubscriptionPayload = {
    player_id: playerStore.player.id,
  }

  websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.GAME_SUBSCRIPTION, payload))
}

const playerPseudo = ref(playerStore.player?.pseudo ?? '')

watch(playerPseudo, () => {
  if (!playerStore.player) return

  const payload: UpdatePlayerPseudoPayload = {
    player_id: playerStore.player.id,
    new_pseudo: playerPseudo.value,
  }

  websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.UPDATE_PLAYER_PSEUDO, payload))

  playerStore.player.pseudo = playerPseudo.value
})

onUnmounted(() => {
  websocketService.unsubscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE)
  websocketService.unsubscribe(WS_MESSAGES_TYPE.WAITING_GAME)
  websocketService.unsubscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME)
})
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-2">
    <h1 class="text-4xl font-bold">Strategia</h1>
    <br />
    <input
      v-model="playerPseudo"
      class="bg-transparent text-sm border border-slate-200 rounded-md px-3 py-2 w-96"
      :placeholder="playerStore.player?.pseudo"
    />
    <br />
    <button
      class="bg-blue-500 text-white px-4 py-2 rounded w-96"
      v-if="!pendingGameStore.isGameLaunching"
      :class="{ 'bg-green-500': pendingGameStore.isSubscribedToGame }"
      @click="sendSubscriptionToGame"
    >
      Join next game <br />
      ({{ pendingGameStore.numberOfWaitingPlayers }} players waiting) <br />
      Launching in {{ pendingGameStore.secondsBeforeLaunch }} seconds
    </button>
    <button
      class="bg-blue-500 text-white px-4 py-2 rounded w-96"
      v-else
      :class="{ 'bg-green-500': pendingGameStore.isSubscribedToGame }"
    >
      Game launching ! <br />
      ({{ pendingGameStore.numberOfWaitingPlayers }} players waiting)
    </button>
  </div>
</template>
