import {defineStore} from "pinia";
import type Player from "@/models/player.ts";

export const usePlayerStore = defineStore('playerStore', {
  state: () => ({
    player: null as Player | null,
  }),
  getters: {
    hasConnectedPlayer(state): boolean {
      return state.player !== null;
    }
  },
  actions: {
    setPlayer(player: Player) {
      this.player = player;
    },
    unsetPlayer() {
      this.player = null;
    }
  }
});
