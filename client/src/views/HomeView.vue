<script setup lang="ts">
import { onUnmounted } from 'vue'
import container from '@/container/container'
import JoinGameComponent from '@/components/JoinGameComponent.vue'
import PlayerPseudoComponent from '@/components/PlayerPseudoComponent.vue'
import HomeViewPresenter from '@/presenter/home-view.presenter'

import { onMounted } from 'vue'
const homeViewPresenter = container.get(HomeViewPresenter)

onMounted(() => {
  homeViewPresenter.handlePageMounted()
})

onUnmounted(() => {
  homeViewPresenter.handlePageUnmounted()
})
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-2">
    <h1 class="text-4xl font-bold">Strategia</h1>
    <br />
    <PlayerPseudoComponent />
    <br />
    <JoinGameComponent />
  </div>
</template>
