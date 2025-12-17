<script setup lang="ts">

import Game from "@/components/Game.vue";
import {useRoute, useRouter} from "vue-router";
import {WebSocketService} from "@/services/websocket.service.ts";
import container from "@/container/container.ts";
import type GameUnsubscribePayload from "@/ws-exchange/game-unsubscribe-payload.ts";
import type {WsExchangeTemplate} from "@/ws-exchange/ws-exchange-template.ts";
import {usePlayerStore} from "@/stores/player.store.ts";
import type ExitGamePayload from "@/ws-exchange/exit-game-payload.ts";

const router = useRouter();
const websocketService: WebSocketService = container.get(WebSocketService);
const playerStore = usePlayerStore();
const route = useRoute()
const gameId = route.query.gameId as string // TODO => la répétition avec le component game est pas dingue

function redirectToHome() {
  websocketService.unsubscribe("server_update_datas");
  const exitGamePayload: ExitGamePayload = {
    player_id: playerStore.player!.id,
    game_id: gameId,
  }

  const exitGameExchange: WsExchangeTemplate<ExitGamePayload> = {
    type: "exit_game",
    payload: exitGamePayload
  }

  websocketService.send<ExitGamePayload>(exitGameExchange);

  router.push('/')
}

</script>

<template>
  <button @click="redirectToHome()">Access Home</button>
  <br />
  <Game />
</template>
