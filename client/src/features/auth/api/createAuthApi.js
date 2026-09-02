import { Session } from "@/entities/session";
import { MESSAGE_TYPES } from "@/shared/api/websocket";
import { clearStoredSession } from "@/shared/lib";

/** @param {import('@/shared/api/websocket/WebSocketTransport.js').WebSocketTransport} transport */
export const createAuthApi = (transport) => ({
    /** @param {import('./contracts.js').CredentialsRequest} credentials */
    register: (credentials) =>
        transport.request(MESSAGE_TYPES.REGISTER, credentials),
    /** @param {import('./contracts.js').CredentialsRequest} credentials */
    login: async (credentials) => {
        const dto = await transport.request(MESSAGE_TYPES.LOGIN, credentials);
        const session = new Session(
      /** @type {import('@/entities/session/api/contracts.js').SessionDto} */(
                dto
            ),
        );
        localStorage.setItem("session_id", session.sessionId);
        return session;
    },
    logout: async () => {
        const sessionId = localStorage.getItem("session_id");
        if (sessionId)
            await transport.request(MESSAGE_TYPES.LOGOUT, { session_id: sessionId });
        clearStoredSession();
    },
});
