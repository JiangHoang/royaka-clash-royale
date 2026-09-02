export class Session {
    /** @param {import('../api/contracts.js').SessionDto} data */
    constructor(data) { this.sessionId = data.session_id; this.expiresAt = data.expires_at; }
}
