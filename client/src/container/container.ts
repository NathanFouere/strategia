import { Container } from 'inversify'
import { WebSocketService } from '@/services/websocket.service.ts'
import { COMMON_DEPENDENCY_TYPES } from '@/container/common.types.ts'
import PlayerPseudoComponentPresenter from '@/presenter/player-pseudo-component.presenter'
import JoinGameComponentPresenter from '@/presenter/join-game-component.presenter'
import HomeViewPresenter from '@/presenter/home-view.presenter'
import GameComponentPresenter from '@/presenter/game-component.presenter'
import LeaveGameComponentPresenter from '@/presenter/leave-game-component.presenter'
import GameViewPresenter from '@/presenter/game-view.presenter'

const container: Container = new Container()

container
  .bind<string>(COMMON_DEPENDENCY_TYPES.WebSocketUrl)
  .toConstantValue('ws://localhost:8080/ws')
container.bind(WebSocketService).toSelf().inSingletonScope()
container.bind(PlayerPseudoComponentPresenter).toSelf()
container.bind(JoinGameComponentPresenter).toSelf()
container.bind(HomeViewPresenter).toSelf()
container.bind(GameComponentPresenter).toSelf()
container.bind(LeaveGameComponentPresenter).toSelf()
container.bind(GameViewPresenter).toSelf()

export default container
