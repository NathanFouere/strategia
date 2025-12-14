<script setup lang="ts">
import { RouterView } from 'vue-router'

import Game from "@/components/Game.vue";
</script>

<template>
  <RouterLink to="/">Access Home</RouterLink>
  <br />
  <Game />
</template>
