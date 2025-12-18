<script setup lang="ts">
import { onMounted } from "vue";
import {usePlayerStore} from "@/stores/player.store.ts";
import {WebSocketService} from "@/services/websocket.service.ts";
import container from "@/container/container.ts";
import { ServerUpdatePayload } from "@/ws-exchange/server-update-payload.ts";
import type PixelClickPayload from "@/ws-exchange/pixel-click-payload.ts";
import {useRoute} from "vue-router";
import type {WsExchangeTemplate} from "@/ws-exchange/ws-exchange-template.ts";

const playerStore = usePlayerStore();
const websocketService: WebSocketService = container.get(WebSocketService);
const route = useRoute()
const gameId = route.query.gameId as string

onMounted(() => {
  const canvas = document.getElementById("grid") as HTMLCanvasElement;
  const ctx = canvas.getContext("2d");

  const cb = (e: ServerUpdatePayload) => {
    console.log("received server_updates_datas", e)
    for (let update_data of e.update_datas) {
      ctx.fillStyle = update_data.color;
      ctx.fillRect(update_data.x, update_data.y, 1, 1);
    }
  };

  websocketService.subscribe<ServerUpdatePayload>("server_update_datas", cb);

  canvas.addEventListener("click", (evt) => {
    console.log("send")
    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const scaleY = canvas.height / rect.height;

    const x = Math.floor((evt.clientX - rect.left) * scaleX);
    const y = Math.floor((evt.clientY - rect.top) * scaleY);

    const pixelClickPayload: PixelClickPayload = {
      id_player: playerStore.player!.id,
      game_id: gameId, // TODO => le recuperer de l'url
      x: x,
      y: y,
    }

    const wsExchange: WsExchangeTemplate<PixelClickPayload> = {
      type: "pixel_click_evt",
      payload: pixelClickPayload,
    }
    websocketService.send<PixelClickPayload>(wsExchange);
  });

  for (let y = 0; y < 1000; y++) {
    for (let x = 0; x < 1000; x++) {
      ctx.fillStyle = "#aaa";
      ctx.fillRect(x, y, 1, 1);
    }
  }
});
</script>

<template>
  <canvas id="grid" width="1000" height="1000"></canvas>
</template>
