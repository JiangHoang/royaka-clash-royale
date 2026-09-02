import { mapBattleEntity } from "@/entities/battle";
import { Player } from "@/entities/player";
import { Tower } from "@/entities/tower";

/** @param {import('./contracts.js').GameResponse} dto */
export const mapGameResponse = (dto) => ({ user: new Player(dto.user), opponent: new Player(dto.opponent), player1: dto.player1 ?? "", map: (dto.map ?? []).map(mapBattleEntity), time: dto.time ?? 0, turn: dto.turn ?? "" });
/** @param {import('./contracts.js').AttackResponse} dto */
export const mapAttackResponse = (dto) => ({ attacker: new Player(dto.attacker), defender: new Player(dto.defender), troop: dto.troop, target: dto.target, damage: dto.damage, isCrit: dto.isCrit, isDestroyed: dto.isDestroyed, turn: dto.turn });
/** @param {import('./contracts.js').HealResponse} dto */
export const mapHealResponse = (dto) => ({ player: new Player(dto.player), opponent: new Player(dto.opponent), troop: dto.troop, healedTower: new Tower(dto.healedTower), healAmount: dto.healAmount, turn: dto.turn });
/** @param {import('./contracts.js').SkipTurnResponse} dto */
export const mapTurnResponse = (dto) => ({ turn: dto.turn, player1: new Player(dto.player1), player2: new Player(dto.player2) });
/** @param {import('./contracts.js').PlayerEvent} dto */
export const mapPlayerEvent = (dto) => ({ player: new Player(dto.player) });
/** @param {import('./contracts.js').GameStateEvent} dto */
export const mapGameStateEvent = (dto) => ({ ...dto, battleMap: dto.battleMap.map(mapBattleEntity) });
/** @param {import('./contracts.js').GameOverEvent} dto */
export const mapGameOverEvent = (dto) => ({ winner: dto.winner ? new Player(dto.winner) : null });
