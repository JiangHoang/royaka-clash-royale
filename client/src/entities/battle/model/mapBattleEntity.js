import { TroopInstance } from "./TroopInstance.js";
import { TowerInstance } from "./TowerInstance.js";
/** @param {import('../api/contracts.js').BattleEntityDto} data */
export const mapBattleEntity = (data) => data.type_entity === "tower" ? new TowerInstance(data) : new TroopInstance(data);
