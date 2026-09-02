export class Troop {
    /** @param {import('../api/contracts.js').TroopDto} data */
    constructor(data) {
        this.name = data.name; this.maxHP = data.max_hp; this.hp = data.hp; this.dmg = data.dmg;
        this.atk = data.atk; this.def = data.def; this.mana = data.mana; this.crit = data.crit;
        this.exp = data.exp; this.speed = data.speed; this.range = data.range; this.type = data.type;
        this.image = data.image; this.description = data.description; this.attackSpeed = data.attack_speed;
        this.aggroPriority = data.aggro_priority; this.rarity = data.rarity;
    }
    healthPercent() { return this.maxHP > 0 ? (this.hp / this.maxHP) * 100 : 0; }
}
