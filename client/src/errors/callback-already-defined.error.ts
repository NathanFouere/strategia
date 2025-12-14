export class CallbackAlreadyDefinedError extends Error {
  constructor(
    public readonly type: string,
  ) {
    super("callback for type : " + type + " is already defined !");
  }
}
