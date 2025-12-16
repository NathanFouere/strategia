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

const playerStore = usePlayerStore();
const pendingGameStore = usePendingGameStore();
const websocketService: WebSocketService = container.get(WebSocketService);

const cb = (e: ConnectionPayload) => {
  const player: Player = {
    id: e.player_id,
    pseudo: e.player_pseudo,
  }

  playerStore.setPlayer(player);
}
websocketService.subscribe<ConnectionPayload>("connexion_exchange", cb)

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
  {{pendingGameStore.isGameLaunching}}
  {{pendingGameStore.isSubscribedToGame}}
  <p>Pending game id: {{pendingGameStore.pendingGameId}}</p>
  <p>Seconds before launch {{pendingGameStore.secondsBeforeLaunch}}</p>
  <button :style="pendingGameStore.isSubscribedToGame ? { 'background-color': 'green'}: ''" @click="sendSubscriptionToGame">
    Join next game {{pendingGameStore.pendingGameId}} ({{pendingGameStore.numberOfWaitingPlayers}} players waiting)
  </button>
  <br />
</template>

<style scoped>

</style>
