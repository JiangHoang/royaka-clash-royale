import { Match } from "@/entities/match";
import { MESSAGE_TYPES } from "@/shared/api/websocket";
/** @param {import('@/shared/api/websocket/WebSocketTransport.js').WebSocketTransport} transport */
export const createMatchApi = (transport) => ({
    /** @param {import('./contracts.js').GameMode} mode */
    find: (mode) => transport.request(MESSAGE_TYPES.FIND_MATCH, { username: localStorage.getItem("username") ?? "", mode }),
    /** @param {(match: Match) => void} handler */
    onFound: (handler) => transport.on(MESSAGE_TYPES.MATCH_FOUND, (raw, response) => { if (response.success) handler(new Match(/** @type {import('@/entities/match/api/contracts.js').MatchDto} */ (raw))); }),
    /** @param {() => void} handler */
    onTimeout: (handler) => transport.on(MESSAGE_TYPES.MATCH_TIMEOUT, (_raw, response) => { if (response.success) handler(); }),
});
