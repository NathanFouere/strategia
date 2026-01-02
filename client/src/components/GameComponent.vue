<script setup lang="ts">
import { onMounted } from 'vue'
import { usePlayerStore } from '@/stores/player.store.ts'
import { WebSocketService } from '@/services/websocket.service.ts'
import container from '@/container/container.ts'
import type PixelClickPayload from '@/ws-exchange/pixel-click-payload.ts'
import { useRoute } from 'vue-router'
import { WS_MESSAGES_TYPE } from '@/ws-exchange/ws-exchange-template.ts'
import GenerateWsTemplate from '@/utils/generate-ws-template.ts'
import type ServerUpdatePayload from '@/ws-exchange/server-update-payload'

const playerStore = usePlayerStore()
const websocketService: WebSocketService = container.get(WebSocketService)
const route = useRoute()
const gameId = route.query.gameId as string

onMounted(() => {
  const canvas = document.getElementById('grid') as HTMLCanvasElement
  const ctx = canvas.getContext('2d')

  const cb = (e: ServerUpdatePayload) => {
    if (ctx == null) {
      return
    }
    for (const update_data of e.update_datas) {
      ctx.fillStyle = update_data.color
      ctx.fillRect(update_data.x, update_data.y, 1, 1)
    }
  }

  websocketService.subscribe(WS_MESSAGES_TYPE.SERVER_UPDATE_DATAS, cb)

  canvas.addEventListener('click', (evt) => {
    const rect = canvas.getBoundingClientRect()
    const scaleX = canvas.width / rect.width
    const scaleY = canvas.height / rect.height

    const x = Math.floor((evt.clientX - rect.left) * scaleX)
    const y = Math.floor((evt.clientY - rect.top) * scaleY)

    const pixelClickPayload: PixelClickPayload = {
      id_player: playerStore.player!.id,
      game_id: gameId,
      x: x,
      y: y,
    }

    websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.PIXEL_CLICK_EVT, pixelClickPayload))
  })

  if (ctx == null) {
    return
  }

  for (let y = 0; y < 1000; y++) {
    for (let x = 0; x < 1000; x++) {
      ctx.fillStyle = '#aaa'
      ctx.fillRect(x, y, 1, 1)
    }
  }
})
</script>

<template>
  <canvas id="grid" width="1000" height="1000"></canvas>
</template>
