<script setup lang="ts">
defineProps<{
  startProgressionPercentage: number
}>()
</script>

<template>
  <div class="w-full bg-gray-200">
    <div
      class="bg-blue-500 h-4 transition-all duration-300"
      :style="{ width: startProgressionPercentage + '%' }"
    ></div>
  </div>
</template>
