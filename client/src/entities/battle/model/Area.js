import { Position } from "./Position.js";
export class Area {
    /** @param {import('../api/contracts.js').AreaDto} data */
    constructor(data) { this.topLeft = new Position(data.top_left); this.bottomRight = new Position(data.bottom_right); }
}
