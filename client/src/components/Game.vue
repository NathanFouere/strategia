<script setup lang="ts">
import { onMounted } from "vue";
import {usePlayerStore} from "@/stores/player.store.ts";
import {WebSocketService} from "@/services/websocket.service.ts";
import container from "@/container/container.ts";
import type {ServerUpdatePayload} from "@/ws-exchange/server-update-payload.ts";
import type ConnectionPayload from "@/ws-exchange/connection-payload.ts";
import type Player from "@/models/player.ts";

const playerStore = usePlayerStore();
const websocketService: WebSocketService = container.get(WebSocketService);


onMounted(() => {
  const canvas = document.getElementById("grid") as HTMLCanvasElement;
  const ctx = canvas.getContext("2d");

  // TODO => montre qu'il y a problematique d'unifier ça
  const cb1 = (e: ConnectionPayload) => {
    const player: Player = {
      id: e.player_id,
      pseudo: e.player_pseudo,
    }
    console.log("im called");
    playerStore.setPlayer(player);
  }
  websocketService.subscribe<ConnectionPayload>("connexion-exchange", cb1)

  const cb = (e: ServerUpdatePayload) => {
    console.log("je reçois un truc")
    for (let update_data of e.update_datas) {
      console.log("updata data ", update_data)
      ctx.fillStyle = update_data.color;
      ctx.fillRect(update_data.x, update_data.y, 1, 1);
    }
  };

  websocketService.subscribe<ServerUpdatePayload>("server-update-datas", cb);

    canvas.addEventListener("click", (evt) => {
      const rect = canvas.getBoundingClientRect();
      const scaleX = canvas.width / rect.width;
      const scaleY = canvas.height / rect.height;

      const x = Math.floor((evt.clientX - rect.left) * scaleX);
      const y = Math.floor((evt.clientY - rect.top) * scaleY);

      const msg = JSON.stringify({ x : x, y : y, id_player: playerStore.player?.id });
      websocketService.ws.send(msg);
    });

  for (let y = 0; y < 100; y++) {
    for (let x = 0; x < 100; x++) {
      ctx.fillStyle = "#aaa";
      ctx.fillRect(x, y, 1, 1);
    }
  }
});
</script>

<template>
  <canvas id="grid" width="100" height="100"></canvas>
</template>
