<script setup lang="ts">
import { useRoute } from 'vue-router'
import GameComponent from '@/components/GameComponent.vue'
import LeaveGameComponent from '@/components/LeaveGameComponent.vue'

const route = useRoute()
const gameId = route.query.gameId as string // TODO => la répétition avec le component game est pas dingue
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-2">
    <GameComponent :game-id="gameId" />
    <br />
    <LeaveGameComponent :game-id="gameId" />
  </div>
</template>
