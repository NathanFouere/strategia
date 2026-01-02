import {
  type ExchangeTypes,
  type Payloads,
  type WsExchangeTemplate,
} from '@/ws-exchange/ws-exchange-template.ts'

export default function GenerateWsTemplate(
  str: ExchangeTypes,
  payload: Payloads,
): WsExchangeTemplate<typeof str> {
  return {
    type: str,
    payload: payload,
  }
}
