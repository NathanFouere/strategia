<script setup lang="ts">
import { onMounted } from 'vue'
import container from '@/container/container.ts'
import { registerServerUpdateHandler } from '@/ws-handler/server-update.handler'
import GameComponentPresenter from '@/presenter/game-component.presenter'

const gameComponentPresenter: GameComponentPresenter = container.get(GameComponentPresenter)
const props = defineProps<{
  gameId: string
}>()
onMounted(() => {
  const canvas = document.getElementById('grid') as HTMLCanvasElement
  const ctx = canvas.getContext('2d')
  if (ctx == null) {
    throw new Error('CTX is undefined')
  }

  registerServerUpdateHandler(ctx)

  canvas.addEventListener('click', (evt: MouseEvent) => {
    const rect = canvas.getBoundingClientRect()
    const scaleX = canvas.width / rect.width
    const scaleY = canvas.height / rect.height

    const x = Math.floor((evt.clientX - rect.left) * scaleX)
    const y = Math.floor((evt.clientY - rect.top) * scaleY)

    gameComponentPresenter.handleClickOnBoard(props.gameId, x, y)
  })

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
