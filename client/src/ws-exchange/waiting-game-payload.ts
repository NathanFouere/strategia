export default interface WaitingGamePayload {
  seconds_before_launch: number;
  game_id: string;
  number_of_waiting_players: number;
  is_player_waiting_for_game: boolean;
  is_game_launching: boolean;
}
