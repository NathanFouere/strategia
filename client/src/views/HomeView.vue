<script setup lang="ts">
import { RouterLink } from 'vue-router'
import {usePlayerStore} from "@/stores/player.store.ts";
import type Player from "@/models/player.ts";
import {WebSocketService} from "@/services/websocket.service.ts";
import type ConnectionPayload from "@/ws-exchange/connection-payload.ts";
import container from "@/container/container.ts";

const playerStore = usePlayerStore();
const websocketService: WebSocketService = container.get(WebSocketService);

const cb = (e: ConnectionPayload) => {
  const player: Player = {
    id: e.player_id,
    pseudo: e.player_pseudo,
  }
  console.log("im called");
  playerStore.setPlayer(player);
}
websocketService.subscribe<ConnectionPayload>("connexion-exchange", cb)
</script>

<template>
  <p>Home</p>
  <br/>
  <RouterLink to="/game">Join next game</RouterLink>
</template>

<style scoped>

</style>
