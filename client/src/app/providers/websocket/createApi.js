import { createAuthApi } from "@/features/auth";
import { createDeckApi } from "@/features/deck";
import { createGameApi } from "@/features/game";
import { createMatchApi } from "@/features/matchmaking";
import { createUserApi } from "@/features/user";

/** @param {import('@/shared/api/websocket/WebSocketTransport.js').WebSocketTransport} transport */
export const createApi = (transport) => ({
    auth: createAuthApi(transport),
    user: createUserApi(transport),
    deck: createDeckApi(transport),
    match: createMatchApi(transport),
    game: createGameApi(transport),
});
