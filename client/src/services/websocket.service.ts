import type {WsExchangeTemplate} from "@/ws-exchange/ws-exchange-template.ts";
import {UndefinedCallbackError} from "@/errors/undefined-callback.error.ts";
import {CallbackAlreadyDefinedError} from "@/errors/callback-already-defined.error.ts";
import {ConnectionNotReadyError} from "@/errors/connection-not-ready.error.ts";
import {inject, injectable} from "inversify";
import {COMMON_DEPENDENCY_TYPES} from "@/container/common.types.ts";

// TODO => se débarasser des any ici, se renseigner sur "discriminated unions" cf . https://www.codefixeshub.com/typescript/discriminated-unions-type-narrowing-for-objects-wi
@injectable()
export class WebSocketService {
  // TODO => enlever le fait que ce soit public
  public ws: WebSocket;
  private callbackDict: Record<string, ((payload: any) => void )> = {};

  constructor(
    @inject(COMMON_DEPENDENCY_TYPES.WebSocketUrl) url: string,
  ) {
    console.log("CONSTRUCTOR CALLED")
    this.ws = new WebSocket(url);
    this.ws.onmessage = (event: MessageEvent): void => {
      const message = JSON.parse(event.data) as WsExchangeTemplate<any>;
      const callback = this.callbackDict[message.type];
      if (callback === undefined) {
        throw new UndefinedCallbackError(message.type)
      }
      callback(message.payload);
    }
  }

  public subscribe<T>(type: string, callback: (payload: T) => void): void {
    if (this.callbackDict[type] !== undefined) {
      throw new CallbackAlreadyDefinedError(type);
    }
    this.callbackDict[type] = callback;
  }

  public send<T>(message: WsExchangeTemplate<T>): void{
    const serializedMessage = JSON.stringify(message);
    this.ws.send(serializedMessage)
  }

  public unsubscribe(type: string): void {
    if (this.callbackDict[type] === undefined) {
      throw new UndefinedCallbackError(type)
    }
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
