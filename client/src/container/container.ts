import { Container } from 'inversify'
import { WebSocketService } from '@/services/websocket.service.ts'
import { COMMON_DEPENDENCY_TYPES } from '@/container/common.types.ts'
import PlayerPseudoComponentPresenter from '@/presenter/player-pseudo-component.presenter'
import JoinGameComponentPresenter from '@/presenter/join-game-component.presenter'
import HomeViewPresenter from '@/presenter/home-view.presenter'
import GameComponentPresenter from '@/presenter/game-component.presenter'
import LeaveGameComponentPresenter from '@/presenter/leave-game-component.presenter'
import GameViewPresenter from '@/presenter/game-view.presenter'
import GameStartupStatusHandler from '@/ws-handler/game-startup-status.handler'
import ConnectionHandler from '@/ws-handler/connection.handler'
import WaitingGameHandler from '@/ws-handler/waiting-game.handler'
import RedirectToGameHandler from '@/ws-handler/redirect-to-game.handler'
import ServerUpdateHandler from '@/ws-handler/server-update.handler'

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
container.bind(GameStartupStatusHandler).toSelf()
container.bind(ConnectionHandler).toSelf()
container.bind(WaitingGameHandler).toSelf()
container.bind(RedirectToGameHandler).toSelf()
container.bind(ServerUpdateHandler).toSelf()

export default container
