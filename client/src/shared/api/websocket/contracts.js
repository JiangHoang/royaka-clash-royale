import { isMessageType } from "./messageTypes.js";

/** @typedef {{ code: string, message: string }} WsError */
/** @template T @typedef {{ type: import('./messageTypes.js').MessageType, request_id: string, data: T }} WsRequest */
/** @template T @typedef {{ type: import('./messageTypes.js').MessageType, request_id?: string, success: boolean, message?: string, error?: WsError | null, data?: T }} WsResponse */
/** @template T @typedef {{ type: import('./messageTypes.js').MessageType, success: boolean, message?: string, error?: WsError | null, data?: T }} WsEvent */

/** @param {unknown} value @returns {value is Record<string, unknown>} */
export const isRecord = (value) => typeof value === "object" && value !== null && !Array.isArray(value);

/** @param {unknown} value @returns {value is WsResponse<unknown>} */
export function isWsResponse(value) {
    if (!isRecord(value)) return false;
    return isMessageType(value.type)
        && typeof value.success === "boolean"
        && (value.request_id === undefined || typeof value.request_id === "string")
        && (value.message === undefined || typeof value.message === "string")
        && (value.error === undefined || value.error === null || isRecord(value.error));
}

export {};
