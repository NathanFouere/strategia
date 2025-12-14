export class UndefinedCallbackError extends Error {
  constructor(
    public readonly type: string,
  ) {
    super("callback for type : " + type + " is not defined !");
  }
}
