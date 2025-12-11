<script setup lang="ts">
import { onMounted } from "vue";
let playerId: string | null = null

const log = (m: any) => {
  (document.getElementById("log")!.textContent += m + "\n");
}

onMounted(() => {
  const canvas = document.getElementById("grid") as HTMLCanvasElement;
  const ctx = canvas.getContext("2d");
  const ws = new WebSocket("ws://localhost:8080/ws");
  ws.onopen = () => log("WebSocket connecté");
  ws.onmessage = (e) => {
    const data = JSON.parse(e.data);
    if (data.type === "connexion-exchange") {
      playerId = JSON.stringify(data.payload["player-id"])
      log("player id is " + playerId)
    }
    if (data.type === "server-update-datas" && data.payload[0]) {
      log("server update datas" + JSON.stringify(data.payload))
      for (let  key of data.payload) {
        ctx.fillStyle = key.color;
        ctx.fillRect(key.x, key.y, 1, 1);
      }
    }
  };

    canvas.addEventListener("click", (evt) => {
      const rect = canvas.getBoundingClientRect();
      const scaleX = canvas.width / rect.width;
      const scaleY = canvas.height / rect.height;

      const x = Math.floor((evt.clientX - rect.left) * scaleX);
      const y = Math.floor((evt.clientY - rect.top) * scaleY);

      const msg = JSON.stringify({ x : x, y : y, id_player: playerId });
      ws.send(msg);
      log("Envoyé: " + msg);
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
  <pre id="log"></pre>
</template>
