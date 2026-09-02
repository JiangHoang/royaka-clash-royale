export class Tower {
    /** @param {import('../api/contracts.js').TowerDto} data */
    constructor(data) {
        this.type = data.type; this.maxHP = data.max_hp; this.hp = data.hp; this.atk = data.atk;
        this.def = data.def; this.crit = data.crit; this.exp = data.exp; this.range = data.range;
        this.attackSpeed = data.attack_speed;
    }
    healthPercent() { return this.maxHP > 0 ? (this.hp / this.maxHP) * 100 : 0; }
}
