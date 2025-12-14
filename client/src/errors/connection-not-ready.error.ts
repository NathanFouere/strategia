export class ConnectionNotReadyError extends Error {
  constructor() {
    super("No connected websocket !");
  }
}
