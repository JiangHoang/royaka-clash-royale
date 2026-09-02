import { MESSAGE_TYPES } from "@/shared/api/websocket";
import { mapAttackResponse, mapGameOverEvent, mapGameResponse, mapGameStateEvent, mapHealResponse, mapPlayerEvent, mapTurnResponse } from "./mappers.js";

const username = () => localStorage.getItem("username") ?? "";
/** @param {import('@/shared/api/websocket/WebSocketTransport.js').WebSocketTransport} transport */
export const createGameApi = (transport) => {
    /** @template T @param {import('@/shared/api/websocket/messageTypes.js').MessageType} type @param {(dto: never) => T} mapper @param {(data: T, response: import('@/shared/api/websocket/contracts.js').WsResponse<unknown>) => void} handler */
    const onMapped = (type, mapper, handler) => transport.on(type, (raw, response) => { if (response.success) handler(mapper(/** @type {never} */ (raw)), response); });
    return {
        /** @param {string} roomId */ get: async (roomId) => mapGameResponse(/** @type {import('./contracts.js').GameResponse} */ (await transport.request(MESSAGE_TYPES.GET_GAME, { room_id: roomId, username: username() }))),
        /** @param {{roomId: string, troop: string, target: string}} input */ attack: async ({ roomId, troop, target }) => mapAttackResponse(/** @type {import('./contracts.js').AttackResponse} */ (await transport.request(MESSAGE_TYPES.ATTACK, { room_id: roomId, username: username(), troop, target }))),
        /** @param {{roomId: string, troop: string}} input */ heal: async ({ roomId, troop }) => mapHealResponse(/** @type {import('./contracts.js').HealResponse} */ (await transport.request(MESSAGE_TYPES.HEAL, { room_id: roomId, username: username(), troop }))),
        /** @param {string} roomId */ skipTurn: async (roomId) => mapTurnResponse(/** @type {import('./contracts.js').SkipTurnResponse} */ (await transport.request(MESSAGE_TYPES.SKIP_TURN, { room_id: roomId, username: username() }))),
        /** @param {{roomId: string, troop: string, x: number, y: number}} input */ selectTroop: async ({ roomId, troop, x, y }) => mapPlayerEvent(/** @type {import('./contracts.js').PlayerEvent} */ (await transport.request(MESSAGE_TYPES.SELECT_TROOP, { room_id: roomId, username: username(), troop, x, y }))),
        /** @param {string} roomId */ leave: (roomId) => transport.request(MESSAGE_TYPES.LEAVE_GAME, { room_id: roomId, username: username() }),
        /** @param {string} roomId */ playAgain: (roomId) => transport.request(MESSAGE_TYPES.PLAY_AGAIN, { room_id: roomId }),
        onAttack: (handler) => onMapped(MESSAGE_TYPES.ATTACK_RESPONSE, mapAttackResponse, handler),
        onHeal: (handler) => onMapped(MESSAGE_TYPES.HEAL_RESPONSE, mapHealResponse, handler),
        onTurnChanged: (handler) => onMapped(MESSAGE_TYPES.SKIP_TURN_RESPONSE, mapTurnResponse, handler),
        onTroopChanged: (handler) => onMapped(MESSAGE_TYPES.TROOP_RESPONSE, mapPlayerEvent, handler),
        onManaChanged: (handler) => onMapped(MESSAGE_TYPES.MANA_UPDATE, mapPlayerEvent, handler),
        onStateChanged: (handler) => onMapped(MESSAGE_TYPES.GAME_STATE, mapGameStateEvent, handler),
        onGameOver: (handler) => onMapped(MESSAGE_TYPES.GAME_OVER_RESPONSE, mapGameOverEvent, handler),
    };
};
