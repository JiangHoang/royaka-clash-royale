import { Troop } from "@/entities/troop";
import { Position } from "./Position.js";
export class TroopInstance {
    /** @param {import('../api/contracts.js').TroopInstanceDto} data */
    constructor(data) {
        this.id = data.id; this.template = new Troop(data.template); this.typeEntity = data.type_entity;
        this.owner = data.owner; this.position = new Position(data.position); this.isDead = data.is_dead;
        this.lastAttack = data.last_attack ? new Date(data.last_attack) : null;
    }
}
