import type ConnectionPayload from '@/ws-exchange/connection-payload.ts'
import type WaitingGamePayload from '@/ws-exchange/waiting-game-payload.ts'
import type RedirectToGamePayload from '@/ws-exchange/redirect-to-game-payload.ts'
import type ServerUpdatePayload from '@/ws-exchange/server-update-payload.ts'
import type ExitGamePayload from '@/ws-exchange/exit-game-payload.ts'
import type SetInWaitingLobbyPayload from '@/ws-exchange/set-in-waiting-lobby-payload.ts'
import type GameUnsubscribePayload from '@/ws-exchange/game-unsubscribe-payload.ts'
import type GameSubscriptionPayload from '@/ws-exchange/game-subscription-payload.ts'
import type UpdatePlayerPseudoPayload from '@/ws-exchange/update-player-pseudo-payload.ts'
import type PixelClickPayload from '@/ws-exchange/pixel-click-payload.ts'

export const WS_MESSAGES_TYPE = {
  CONNECTION_EXCHANGE: 'connexion_exchange',
  WAITING_FOR_GAME_EXCHANGE: 'waiting_game_exchange',
  REDIRECT_TO_GAME: 'redirect_to_game',
  SERVER_UPDATE_DATAS: 'server_update_datas',
  EXIT_GAME: 'exit_game',
  SET_IN_WAITING_LOBBY: 'set_in_waiting_lobby',
  WAITING_GAME: 'waiting_game_exchange',
  GAME_UNSUBSCRIBE: 'game_unsubscribe',
  GAME_SUBSCRIPTION: 'game_subscription',
  UPDATE_PLAYER_PSEUDO: 'update_player_pseudo',
  PIXEL_CLICK_EVT: 'pixel_click_evt',
} as const

// cf computed property name . https://dev.to/ahmad_tibibi/ts1166-a-computed-property-name-in-a-class-property-declaration-must-have-a-simple-literal-type-3k9j
export type WsPayloadPerType = {
  [WS_MESSAGES_TYPE.CONNECTION_EXCHANGE]: ConnectionPayload
  [WS_MESSAGES_TYPE.WAITING_FOR_GAME_EXCHANGE]: WaitingGamePayload
  [WS_MESSAGES_TYPE.REDIRECT_TO_GAME]: RedirectToGamePayload
  [WS_MESSAGES_TYPE.SERVER_UPDATE_DATAS]: ServerUpdatePayload
  [WS_MESSAGES_TYPE.EXIT_GAME]: ExitGamePayload
  [WS_MESSAGES_TYPE.SET_IN_WAITING_LOBBY]: SetInWaitingLobbyPayload
  [WS_MESSAGES_TYPE.WAITING_GAME]: WaitingGamePayload
  [WS_MESSAGES_TYPE.GAME_UNSUBSCRIBE]: GameUnsubscribePayload
  [WS_MESSAGES_TYPE.GAME_SUBSCRIPTION]: GameSubscriptionPayload
  [WS_MESSAGES_TYPE.UPDATE_PLAYER_PSEUDO]: UpdatePlayerPseudoPayload
  [WS_MESSAGES_TYPE.PIXEL_CLICK_EVT]: PixelClickPayload
}

export interface WsExchangeTemplate<T extends ExchangeTypes> {
  type: T
  payload: WsPayloadPerType[T]
}

// permet de générer un objet de types
export type WsExchangeMessage = {
  [K in keyof WsPayloadPerType]: {
    type: K
    payload: WsPayloadPerType[K]
  }
}
// permet de générer une union de type, ici les clés de WsPayloadPerType
export type ExchangeTypes = keyof WsPayloadPerType

// permet de générer une union de type, ici les valeur de WsPayloadPerType
export type Payloads = WsPayloadPerType[keyof WsPayloadPerType]
