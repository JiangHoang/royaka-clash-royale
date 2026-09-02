import { Troop } from "@/entities/troop";
import { MESSAGE_TYPES } from "@/shared/api/websocket";
/** @param {import('@/shared/api/websocket/WebSocketTransport.js').WebSocketTransport} transport */
export const createDeckApi = (transport) => ({
    get: async () => /** @type {import('./contracts.js').DeckResponse} */ (await transport.request(MESSAGE_TYPES.GET_DECK, {})).map((dto) => new Troop(dto)),
});
