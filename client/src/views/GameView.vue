<script setup lang="ts">
import { useRoute } from 'vue-router'
import GameComponent from '@/components/GameComponent.vue'
import LeaveGameComponent from '@/components/LeaveGameComponent.vue'
import GameStartupProgressBarComponent from '@/components/GameStartupProgressBarComponent.vue'
import container from '@/container/container'
import GameViewPresenter from '@/presenter/game-view.presenter'
import { onMounted } from 'vue'

const route = useRoute()
const gameId = route.query.gameId as string // TODO => la répétition avec le component game est pas dingue

const gameViewPresenter = container.get(GameViewPresenter)

onMounted(() => {
  gameViewPresenter.initialize()
})
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-2">
    <GameComponent :game-id="gameId" />
    <br />
    <LeaveGameComponent :game-id="gameId" />
    <br />
    <GameStartupProgressBarComponent
      v-if="
        !gameViewPresenter.ongoingGameStore.gameStarted &&
        gameViewPresenter.ongoingGameStore.startProgressionPercentage
      "
      :start-progression-percentage="gameViewPresenter.ongoingGameStore.startProgressionPercentage"
    />
  </div>
</template>
