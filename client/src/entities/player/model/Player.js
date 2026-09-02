import { Tower } from "@/entities/tower";
import { Troop } from "@/entities/troop";
import { TowerInstance, TroopInstance } from "@/entities/battle";
import { User } from "@/entities/user";
export class Player {
    /** @param {import('../api/contracts.js').PlayerDto} data */
    constructor(data) {
        this.mana = data.mana;
        this.towers = Object.fromEntries(Object.entries(data.towers).map(([key, value]) => [key, new Tower(value)]));
        this.towerInstances = (data.tower_instances ?? []).map((value) => new TowerInstance(value));
        this.troops = data.troops.map((value) => new Troop(value));
        this.troopInstances = (data.troop_instances ?? []).map((value) => new TroopInstance(value));
        this.active = data.active; this.user = new User(data.user); this.turn = data.turn; this.gold = data.gold;
    }
}
