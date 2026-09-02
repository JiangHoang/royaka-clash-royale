export class User {
    /** @param {import('../api/contracts.js').UserDto} data */
    constructor(data) {
        this.id = data.id;
        this.username = data.username;
        this.createdAt = data.createdAt ? new Date(data.createdAt) : null;
        this.lastLogin = data.lastLogin ? new Date(data.lastLogin) : null;
        this.isActive = Boolean(data.isActive);
        this.exp = data.exp;
        this.level = data.level;
        this.gamesPlayed = data.gamesPlayed;
        this.gamesWon = data.gamesWon;
        this.avatar = data.avatar || "1";
        this.gold = data.gold;
    }
}
