import { Container } from 'inversify'
import { WebSocketService } from '@/services/websocket.service.ts'
import { COMMON_DEPENDENCY_TYPES } from '@/container/common.types.ts'

const container: Container = new Container()

container
  .bind<string>(COMMON_DEPENDENCY_TYPES.WebSocketUrl)
  .toConstantValue('ws://localhost:8080/ws')
container.bind(WebSocketService).toSelf().inSingletonScope()

export default container
