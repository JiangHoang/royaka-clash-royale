import { Tower } from "@/entities/tower";
import { Area } from "./Area.js";
export class TowerInstance {
    /** @param {import('../api/contracts.js').TowerInstanceDto} data */
    constructor(data) {
        this.id = data.id; this.template = new Tower(data.template); this.typeEntity = data.type_entity;
        this.owner = data.owner; this.area = new Area(data.area); this.isDestroyed = data.is_destroyed;
        this.lastAttack = data.last_attack ? new Date(data.last_attack) : null;
    }
}
