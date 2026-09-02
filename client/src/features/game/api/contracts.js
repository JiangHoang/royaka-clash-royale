/** @typedef {{ room_id: string, username: string }} GameRequest */
/** @typedef {{ room_id: string, username: string, troop: string, target: string }} AttackRequest */
/** @typedef {{ room_id: string, username: string, troop: string }} HealRequest */
/** @typedef {{ room_id: string, username: string, troop: string, x: number, y: number }} SelectTroopRequest */
/** @typedef {{ room_id: string }} PlayAgainRequest */
/** @typedef {{ user: import('@/entities/player/api/contracts.js').PlayerDto, opponent: import('@/entities/player/api/contracts.js').PlayerDto, player1?: string, map?: import('@/entities/battle/api/contracts.js').BattleEntityDto[], time?: number, turn?: string }} GameResponse */
/** @typedef {{ attacker: import('@/entities/player/api/contracts.js').PlayerDto, defender: import('@/entities/player/api/contracts.js').PlayerDto, troop: string, target: string, damage: number, isCrit: boolean, isDestroyed: boolean, turn: string }} AttackResponse */
/** @typedef {{ player: import('@/entities/player/api/contracts.js').PlayerDto, opponent: import('@/entities/player/api/contracts.js').PlayerDto, troop: string, healedTower: import('@/entities/tower/api/contracts.js').TowerDto, healAmount: number, turn: string }} HealResponse */
/** @typedef {{ turn: string, player1: import('@/entities/player/api/contracts.js').PlayerDto, player2: import('@/entities/player/api/contracts.js').PlayerDto }} SkipTurnResponse */
/** @typedef {{ player: import('@/entities/player/api/contracts.js').PlayerDto }} PlayerEvent */
/** @typedef {{ battleMap: import('@/entities/battle/api/contracts.js').BattleEntityDto[], timeLeft: number, player1Guard1: number, player1Guard2: number, player2Guard1: number, player2Guard2: number }} GameStateEvent */
/** @typedef {{ winner?: import('@/entities/player/api/contracts.js').PlayerDto }} GameOverEvent */
export {};
