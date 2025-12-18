<script setup lang="ts">
import { RouterLink } from 'vue-router'
import {usePlayerStore} from "@/stores/player.store.ts";
import type Player from "@/models/player.ts";
import {WebSocketService} from "@/services/websocket.service.ts";
import type ConnectionPayload from "@/ws-exchange/connection-payload.ts";
import container from "@/container/container.ts";
import type WaitingGamePayload from "@/ws-exchange/waiting-game-payload.ts";
import {usePendingGameStore} from "@/stores/pending-game.store.ts";
import type GameSubscriptionPayload from "@/ws-exchange/game-subscription-payload.ts";
import type {WsExchangeTemplate} from "@/ws-exchange/ws-exchange-template.ts";
import router from "@/router";
import type RedirectToGamePayload from "@/ws-exchange/redirect-to-game-payload.ts";
import type GameUnsubscribePayload from "@/ws-exchange/game-unsubscribe-payload.ts";
import type SetInWaitingLobbyPayload from "@/ws-exchange/set-in-waiting-lobby-payload.ts";

const playerStore = usePlayerStore();
const pendingGameStore = usePendingGameStore();
const websocketService: WebSocketService = container.get(WebSocketService);

if (!playerStore.player) {
  const cb = (e: ConnectionPayload) => {
    const player: Player = {
      id: e.player_id,
      pseudo: e.player_pseudo,
    }

    playerStore.setPlayer(player);
  }
  websocketService.subscribe<ConnectionPayload>("connexion_exchange", cb)
} else {
  const setInWaitingLobbyPayload: SetInWaitingLobbyPayload = {
    player_id: playerStore.player!.id
  }

  const setInWaitingLobbyExchange: WsExchangeTemplate<SetInWaitingLobbyPayload> = {
    type: "set_in_waiting_lobby",
    payload: setInWaitingLobbyPayload
  }

  websocketService.send<SetInWaitingLobbyPayload>(setInWaitingLobbyExchange);
}

const cb2 = (e: WaitingGamePayload) => {
  pendingGameStore.setPendingGameId(e.game_id);
  pendingGameStore.setSecondsBeforeLaunch(e.seconds_before_launch);
  pendingGameStore.setNumberOfWaitingPlayers(e.number_of_waiting_players);
  pendingGameStore.setGameLaunching(e.is_game_launching);
  pendingGameStore.setSubscribedToGame(e.is_player_waiting_for_game);
}

websocketService.subscribe<WaitingGamePayload>("waiting_game_exchange", cb2)

const cb3 = (e: RedirectToGamePayload) => {
  websocketService.unsubscribe("connexion_exchange");
  websocketService.unsubscribe("waiting_game_exchange");
  websocketService.unsubscribe("redirect_to_game");
  pendingGameStore.unsetAll();
  router.push('/game?gameId=' + e.game_id);
}

websocketService.subscribe<RedirectToGamePayload>("redirect_to_game", cb3)

// TODO => clarifier le fait que ça inscrive ET desinscrive
function sendSubscriptionToGame(): void {
  if (!playerStore.hasConnectedPlayer) {
    throw new Error("Should have a connected player");
  }
  if (pendingGameStore.isSubscribedToGame) {
    const gameUnsubscribePayload: GameUnsubscribePayload = {
      player_id: playerStore.player!.id
    }

    const gameUnsubscribeExchange: WsExchangeTemplate<GameUnsubscribePayload> = {
      type: "game_unsubscribe",
      payload: gameUnsubscribePayload
    }

    websocketService.send<GameUnsubscribePayload>(gameUnsubscribeExchange);

    return;
  }

  const gameSubscriptionPayload: GameSubscriptionPayload = {
    player_id: playerStore.player!.id
  }

  const gameSubscriptionExchange: WsExchangeTemplate<GameSubscriptionPayload> = {
    type: "game_subscription",
    payload: gameSubscriptionPayload,
  }

  websocketService.send<GameSubscriptionPayload>(gameSubscriptionExchange);
}
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-2">
    <h1 class="text-4xl font-bold">Strategia</h1>

    <br>

    <input class="bg-transparent text-sm border border-slate-200 rounded-md px-3 py-2 w-96" :placeholder="playerStore.player?.id">
    <br />

    <button
      class="bg-blue-500 text-white px-4 py-2 rounded w-96"
      v-if="!pendingGameStore.isGameLaunching"
      :class="{ 'bg-green-500': pendingGameStore.isSubscribedToGame }"
      @click="sendSubscriptionToGame"
    >
      Join next game
      <br />
      ({{pendingGameStore.numberOfWaitingPlayers}} players waiting)
      <br />
      Launching in {{pendingGameStore.secondsBeforeLaunch}} seconds
    </button>


    <button
      class="bg-blue-500 text-white px-4 py-2 rounded w-96"
      v-else
      :class="{ 'bg-green-500': pendingGameStore.isSubscribedToGame }"
    >
      Game launching !
      <br />
      ({{pendingGameStore.numberOfWaitingPlayers}} players waiting)
    </button>
  </div>
</template>


<style scoped>

</style>
