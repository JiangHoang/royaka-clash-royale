/** @typedef {{ x: number, y: number }} PositionDto */
/** @typedef {{ top_left: PositionDto, bottom_right: PositionDto }} AreaDto */
/** @typedef {{ id: string, template: import('../../troop/api/contracts.js').TroopDto, type_entity: 'troop', owner: string, position: PositionDto, is_dead: boolean, last_attack: string }} TroopInstanceDto */
/** @typedef {{ id: string, template: import('../../tower/api/contracts.js').TowerDto, type_entity: 'tower', owner: string, area: AreaDto, is_destroyed: boolean, last_attack: string }} TowerInstanceDto */
/** @typedef {TroopInstanceDto | TowerInstanceDto} BattleEntityDto */
export {};
