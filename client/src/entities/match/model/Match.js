export class Match {
    /** @param {import('../api/contracts.js').MatchDto} data */
    constructor(data) { this.roomId = data.room_id; this.opponent = data.opponent; }
}
