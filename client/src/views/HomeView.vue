<script setup lang="ts">
import { RouterLink } from 'vue-router'
import {usePlayerStore} from "@/stores/player.store.ts";
import type Player from "@/models/player.ts";
import {WebSocketService} from "@/services/websocket.service.ts";
import type ConnectionPayload from "@/ws-exchange/connection-payload.ts";
import container from "@/container/container.ts";
import type WaitingGamePayload from "@/ws-exchange/waiting-game-payload.ts";
import {usePendingGameStore} from "@/stores/pending-game.store.ts";
import type GameSubscriptionPayload from "@/ws-exchange/game-subscription-payload.ts";
import {
  WS_MESSAGES_TYPE,
  type WsExchangeTemplate,
  type WsPayloadPerType
} from "@/ws-exchange/ws-exchange-template.ts";
import router from "@/router";
import type RedirectToGamePayload from "@/ws-exchange/redirect-to-game-payload.ts";
import type GameUnsubscribePayload from "@/ws-exchange/game-unsubscribe-payload.ts";
import type SetInWaitingLobbyPayload from "@/ws-exchange/set-in-waiting-lobby-payload.ts";
import {ref, watch} from "vue";
import type UpdatePlayerPseudoPayload from "@/ws-exchange/update-player-pseudo-payload.ts";
import GenerateWsTemplate from "@/utils/generate-ws-template.ts";

const playerStore = usePlayerStore();
const pendingGameStore = usePendingGameStore();
const websocketService: WebSocketService = container.get(WebSocketService);

if (!playerStore.player) {
  const cb = (payload: ConnectionPayload) => {
    const player: Player = {
      id: payload.player_id,
      pseudo: payload.player_pseudo,
    };

    playerStore.setPlayer(player);
  }
  websocketService.subscribe(WS_MESSAGES_TYPE.CONNECTION_EXCHANGE, cb);
} else {
  const setInWaitingLobbyPayload: SetInWaitingLobbyPayload = {
    player_id: playerStore.player!.id
  };

  websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.SET_IN_WAITING_LOBBY, setInWaitingLobbyPayload));
}

const cb2 = (e: WaitingGamePayload) => {
  pendingGameStore.setPendingGameId(e.game_id);
  pendingGameStore.setSecondsBeforeLaunch(e.seconds_before_launch);
  pendingGameStore.setNumberOfWaitingPlayers(e.number_of_waiting_players);
  pendingGameStore.setGameLaunching(e.is_game_launching);
  pendingGameStore.setSubscribedToGame(e.is_player_waiting_for_game);
}

websocketService.subscribe(WS_MESSAGES_TYPE.WAITING_GAME, cb2)

const cb3 = (e: RedirectToGamePayload) => {
  websocketService.unsubscribe("connexion_exchange");
  websocketService.unsubscribe("waiting_game_exchange");
  websocketService.unsubscribe("redirect_to_game");
  pendingGameStore.unsetAll();
  router.push('/game?gameId=' + e.game_id);
}

websocketService.subscribe(WS_MESSAGES_TYPE.REDIRECT_TO_GAME, cb3)

function sendSubscriptionToGame(): void {
  if (!playerStore.hasConnectedPlayer) {
    throw new Error("Should have a connected player");
  }
  if (pendingGameStore.isSubscribedToGame) {
    const gameUnsubscribePayload: GameUnsubscribePayload = {
      player_id: playerStore.player!.id
    }

    websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.GAME_UNSUBSCRIBE, gameUnsubscribePayload));

    return;
  }

  const gameSubscriptionPayload: GameSubscriptionPayload = {
    player_id: playerStore.player!.id
  }

  websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.GAME_SUBSCRIPTION, gameSubscriptionPayload));
}

let playerPseudo = ref(playerStore.player?.pseudo ?? "");

watch(playerPseudo, () => {{
  if (!playerStore.player) {
    return;
  }

  const updatePlayerPseudoPayload: UpdatePlayerPseudoPayload = {
    player_id: playerStore.player?.id,
    new_pseudo: playerPseudo.value
  }

  websocketService.send(GenerateWsTemplate(WS_MESSAGES_TYPE.UPDATE_PLAYER_PSEUDO, updatePlayerPseudoPayload));
  playerStore.player.pseudo = playerPseudo.value
}})


</script>

<template>
  <div class="flex flex-col items-center gap-4 p-2">
    <h1 class="text-4xl font-bold">Strategia</h1>

    <br>
    <input v-model="playerPseudo" class="bg-transparent text-sm border border-slate-200 rounded-md px-3 py-2 w-96" :placeholder="playerStore.player?.pseudo">
    <br />

    <button
      class="bg-blue-500 text-white px-4 py-2 rounded w-96"
      v-if="!pendingGameStore.isGameLaunching"
      :class="{ 'bg-green-500': pendingGameStore.isSubscribedToGame }"
      @click="sendSubscriptionToGame"
    >
      Join next game
      <br />
      ({{pendingGameStore.numberOfWaitingPlayers}} players waiting)
      <br />
      Launching in {{pendingGameStore.secondsBeforeLaunch}} seconds
    </button>


    <button
      class="bg-blue-500 text-white px-4 py-2 rounded w-96"
      v-else
      :class="{ 'bg-green-500': pendingGameStore.isSubscribedToGame }"
    >
      Game launching !
      <br />
      ({{pendingGameStore.numberOfWaitingPlayers}} players waiting)
    </button>
  </div>
</template>


<style scoped>

</style>
