import type {
  ExchangeTypes,
  Payloads,
  WsExchangeMessage,
  WsExchangeTemplate, WsPayloadPerType
} from "@/ws-exchange/ws-exchange-template.ts";
import {UndefinedCallbackError} from "@/errors/undefined-callback.error.ts";
import {CallbackAlreadyDefinedError} from "@/errors/callback-already-defined.error.ts";
import {ConnectionNotReadyError} from "@/errors/connection-not-ready.error.ts";
import {inject, injectable} from "inversify";
import {COMMON_DEPENDENCY_TYPES} from "@/container/common.types.ts";

@injectable()
export class WebSocketService {
  private ws: WebSocket;
  // Utiliser un any est pas parfait ici, pareil pour le string
  // compliqué à faire autrement et pas vraiment problématique car toutes les méthode publiques du service sont type-safe
  // on ne devrait donc pas alimenter le record avec des choses "fausses"
  // TODO : à corriger tout de même
  private callbackDict: Record<string, ((payload: any) => void )> = {};

  constructor(
    @inject(COMMON_DEPENDENCY_TYPES.WebSocketUrl) url: string,
  ) {
    this.ws = new WebSocket(url);
    this.ws.onmessage = (event: MessageEvent): void => {
      const message = JSON.parse(event.data) as WsExchangeMessage[keyof WsPayloadPerType];
      const callback = this.callbackDict[message.type];
      if (callback === undefined) {
        throw new UndefinedCallbackError(message.type)
      }
      callback(message.payload);
    }
  }

  public subscribe<T extends ExchangeTypes>(type: T, callback: (payload: WsPayloadPerType[T]) => void): void {
    if (this.callbackDict[type] !== undefined) {
      throw new CallbackAlreadyDefinedError(type);
    }
    this.callbackDict[type] = callback;
  }

  public send<T extends ExchangeTypes>(message: WsExchangeTemplate<T>): void{
    const serializedMessage = JSON.stringify(message);
    this.ws.send(serializedMessage)
  }

  public unsubscribe(type: string): void {
    delete this.callbackDict[type] ;
  }

  private checkReadyState(): boolean {
    return this.ws.readyState === 1;
  }

  private assertReadyState(): void {
    if (this.checkReadyState()) {
      return;
    }

    throw new ConnectionNotReadyError();
  }
}
