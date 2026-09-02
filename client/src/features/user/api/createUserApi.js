import { User } from "@/entities/user";
import { MESSAGE_TYPES } from "@/shared/api/websocket";
/** @param {import('@/shared/api/websocket/WebSocketTransport.js').WebSocketTransport} transport */
export const createUserApi = (transport) => ({
    getCurrent: async () => {
        const raw = await transport.request(MESSAGE_TYPES.GET_USER, {});
        const dto = /** @type {import('./contracts.js').CurrentUserResponse} */ (raw);
        return { user: new User(dto.user), maxExp: dto.maxExp };
    },
});
