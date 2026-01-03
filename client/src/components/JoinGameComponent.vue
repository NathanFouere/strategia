<script setup lang="ts">
import container from '@/container/container'
import JoinGameComponentPresenter from '@/presenter/join-game-component.presenter'

const joinGameComponentPresenter = container.get(JoinGameComponentPresenter)
</script>

<template>
  <button
    class="bg-blue-500 text-white px-4 py-2 rounded w-96"
    v-if="!joinGameComponentPresenter.pendingGameStore.isGameLaunching"
    :class="{ 'bg-green-500': joinGameComponentPresenter.pendingGameStore.isSubscribedToGame }"
    @click="joinGameComponentPresenter.sendSubscriptionToGame()"
  >
    Join next game <br />
    ({{ joinGameComponentPresenter.pendingGameStore.numberOfWaitingPlayers }} players waiting)
    <br />
    Launching in {{ joinGameComponentPresenter.pendingGameStore.secondsBeforeLaunch }} seconds
  </button>
  <button
    class="bg-blue-500 text-white px-4 py-2 rounded w-96"
    v-else
    :class="{ 'bg-green-500': joinGameComponentPresenter.pendingGameStore.isSubscribedToGame }"
  >
    Game launching ! <br />
    ({{ joinGameComponentPresenter.pendingGameStore.numberOfWaitingPlayers }} players waiting)
  </button>
</template>
