<script setup lang="ts">
import { ref, watch } from 'vue'
import container from '@/container/container'
import PlayerPseudoComponentPresenter from '@/presenter/player-pseudo-component.presenter'

const playerPseudoComponentPresenter = container.get(PlayerPseudoComponentPresenter)
const playerPseudo = ref(playerPseudoComponentPresenter.playerStore.player?.pseudo ?? '')

watch(playerPseudo, () => {
  if (!playerPseudoComponentPresenter.playerStore.player) return

  playerPseudoComponentPresenter.handlePlayerPseudoChange(playerPseudo.value)
})
</script>

<template>
  <input
    v-model="playerPseudo"
    class="bg-transparent text-sm border border-slate-200 rounded-md px-3 py-2 w-96"
    :placeholder="playerPseudoComponentPresenter.playerStore.player?.pseudo"
  />
</template>
