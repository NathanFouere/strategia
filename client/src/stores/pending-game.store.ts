import {defineStore} from "pinia";

export const usePendingGameStore = defineStore('pendingGameStore', {
  state: () => ({
    pendingGameId: null as string | null,
    secondsBeforeLaunch: null as number | null,
    numberOfWaitingPlayers: null as number | null,
    isSubscribedToGame: false as boolean,
  }),
  actions: {
    setPendingGameId(pendingGameId: string): void {
      this.pendingGameId = pendingGameId;
    },
    setSecondsBeforeLaunch(secondsBeforeLaunch: number): void {
      this.secondsBeforeLaunch = secondsBeforeLaunch;
    },
    setNumberOfWaitingPlayers(numberOfWaitingPlayers: number): void {
      this.numberOfWaitingPlayers = numberOfWaitingPlayers;
    },
    setSubscribedToGame(subscribed: boolean) {
      this.isSubscribedToGame = subscribed;
    },
  }
});
