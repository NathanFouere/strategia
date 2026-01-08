import { defineStore } from 'pinia'

export const useOngoingGameStore = defineStore('ongoingGameStore', {
  state: () => ({
    gameStarted: false as boolean,
    startProgressionPercentage: null as number | null,
  }),

  actions: {
    setGameStarted(started: boolean): void {
      this.gameStarted = started
    },
    setProgressionPercentage(progressionPercentage: number): void {
      this.startProgressionPercentage = progressionPercentage
    },
  },
})
