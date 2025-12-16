import {defineStore} from "pinia";

export const usePendingGameStore = defineStore('pendingGameStore', {
  state: () => ({
    pendingGameId: null as string | null,
    secondsBeforeLaunch: null as number | null,
    numberOfWaitingPlayers: null as number | null,
    isSubscribedToGame: false as boolean,
    isGameLaunching: false as boolean,
  }),
  actions: {
    setPendingGameId(pendingGameId: string): void {
      this.pendingGameId = pendingGameId;
    },
    unsetPendingGameId(): void {
      this.pendingGameId = null
    },
    setSecondsBeforeLaunch(secondsBeforeLaunch: number): void {
      this.secondsBeforeLaunch = secondsBeforeLaunch;
    },
    unsetSecondsBeforeLaunch(): void {
      this.secondsBeforeLaunch = null;
    },
    setNumberOfWaitingPlayers(numberOfWaitingPlayers: number): void {
      this.numberOfWaitingPlayers = numberOfWaitingPlayers;
    },
    unsetNumberOfWaitingsPlayers(): void {
      this.numberOfWaitingPlayers = null;
    },
    setSubscribedToGame(subscribed: boolean): void {
      this.isSubscribedToGame = subscribed;
    },
    unsetSubscribedToGame(): void {
      this.isSubscribedToGame = false
    },
    setGameLaunching(launching: boolean): void {
      this.isGameLaunching = launching;
    },
    unsetGameLaunching(): void {
      this.isGameLaunching = false;
    },
    unsetAll(): void {
      this.unsetNumberOfWaitingsPlayers();
      this.unsetPendingGameId();
      this.unsetSubscribedToGame();
      this.unsetGameLaunching();
      this.unsetSecondsBeforeLaunch();
    },
  }
});
