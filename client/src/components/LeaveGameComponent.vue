<script setup lang="ts">
import container from '@/container/container'
import LeaveGameComponentPresenter from '@/presenter/leave-game-component.presenter'
const leaveGameComponentPresenter: LeaveGameComponentPresenter = container.get(
  LeaveGameComponentPresenter,
)

defineProps<{
  gameId: string
}>()
</script>

<template>
  <button
    class="bg-blue-500 text-white px-4 py-2 rounded"
    @click="leaveGameComponentPresenter.leaveGame(gameId)"
  >
    Leave Game
  </button>
</template>
