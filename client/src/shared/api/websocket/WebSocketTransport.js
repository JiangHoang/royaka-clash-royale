
import { clearStoredSession } from "@/shared/lib/session";
import { isWsResponse } from "./contracts.js";
import { MESSAGE_TYPES } from "./messageTypes.js";
import { WebSocketRequestError } from "./WebSocketRequestError.js";

const REQUEST_TIMEOUT_MS = 10000;

/** @typedef {(data: unknown, response: import('./contracts.js').WsResponse<unknown>) => void} EventListener */
/** @typedef {{ resolve: (value: unknown) => void, reject: (reason: WebSocketRequestError) => void, timeout: number }} PendingRequest */
/** @typedef {{ onStateChange?: (state: 'connecting'|'connected'|'disconnected'|'error') => void, onAuthenticated?: (authenticated: boolean) => void }} TransportCallbacks */

export class WebSocketTransport {
    /** @param {string} url @param {TransportCallbacks} [callbacks] */
    constructor(url, callbacks = {}) {
        this.url = url;
        this.callbacks = callbacks;
        /** @type {WebSocket | null} */ this.socket = null;
        /** @type {Map<string, Set<EventListener>>} */ this.listeners = new Map();
        /** @type {Map<string, PendingRequest>} */ this.pendingRequests = new Map();
        /** @type {number | null} */ this.reconnectTimer = null;
        this.reconnectAttempt = 0;
        this.disposed = false;
    }

    connect() {
        if (!this.url) { this.callbacks.onStateChange?.("error"); return; }
        if (this.socket?.readyState === WebSocket.OPEN || this.socket?.readyState === WebSocket.CONNECTING) return;
        this.callbacks.onStateChange?.("connecting");
        const socket = new WebSocket(this.url);
        this.socket = socket;
        socket.onopen = async () => {
            if (this.socket !== socket) return;
            this.reconnectAttempt = 0;
            this.callbacks.onStateChange?.("connected");
            const sessionId = localStorage.getItem("session_id");
            if (!sessionId) return;
            try {
                await this.request(MESSAGE_TYPES.AUTHENTICATE, { session_id: sessionId });
                this.callbacks.onAuthenticated?.(true);
            } catch {
                clearStoredSession();
                this.callbacks.onAuthenticated?.(false);
            }
        };
        socket.onmessage = (event) => this.handleMessage(event.data);
        socket.onerror = () => this.callbacks.onStateChange?.("error");
        socket.onclose = () => {
            if (this.socket !== socket) return;
            this.socket = null;
            this.rejectPending(new WebSocketRequestError("Connection closed", "connection_closed"));
            this.callbacks.onStateChange?.("disconnected");
            this.callbacks.onAuthenticated?.(false);
            if (!this.disposed) this.scheduleReconnect();
        };
    }

    /** @template TRequest @template TResponse @param {import('./messageTypes.js').MessageType} type @param {TRequest} data @param {number} [timeoutMs] @returns {Promise<TResponse>} */
    request(type, data, timeoutMs = REQUEST_TIMEOUT_MS) {
        if (this.socket?.readyState !== WebSocket.OPEN) return Promise.reject(new WebSocketRequestError("Not connected to server", "not_connected"));
        const requestId = globalThis.crypto?.randomUUID?.() ?? `${Date.now()}-${Math.random()}`;
        return new Promise((resolve, reject) => {
            const timeout = window.setTimeout(() => {
                this.pendingRequests.delete(requestId);
                reject(new WebSocketRequestError(`Request ${type} timed out`, "request_timeout"));
            }, timeoutMs);
            this.pendingRequests.set(requestId, { resolve, reject, timeout });
            /** @type {import('./contracts.js').WsRequest<TRequest>} */
            const request = { type, request_id: requestId, data };
            this.socket?.send(JSON.stringify(request));
        });
    }

    /** @param {import('./messageTypes.js').MessageType} type @param {EventListener} listener */
    on(type, listener) {
        const listeners = this.listeners.get(type) ?? new Set();
        listeners.add(listener);
        this.listeners.set(type, listeners);
        return () => { listeners.delete(listener); if (listeners.size === 0) this.listeners.delete(type); };
    }

    close() {
        this.disposed = true;
        if (this.reconnectTimer !== null) window.clearTimeout(this.reconnectTimer);
        this.rejectPending(new WebSocketRequestError("Transport disposed", "transport_disposed"));
        const socket = this.socket;
        this.socket = null;
        socket?.close();
    }

    /** @param {unknown} rawMessage */
    handleMessage(rawMessage) {
        if (typeof rawMessage !== "string") return;
        /** @type {unknown} */ let parsed;
        try { parsed = JSON.parse(rawMessage); } catch { return; }
        if (!isWsResponse(parsed)) return;
        const requestId = parsed.request_id ?? "";
        const pending = requestId ? this.pendingRequests.get(requestId) : undefined;
        if (pending) {
            window.clearTimeout(pending.timeout);
            this.pendingRequests.delete(requestId);
            if (parsed.success) pending.resolve(parsed.data);
            else pending.reject(new WebSocketRequestError(parsed.error?.message ?? parsed.message ?? "Request failed", parsed.error?.code, parsed));
        }
        this.listeners.get(parsed.type)?.forEach((listener) => listener(parsed.data, parsed));
    }

    /** @param {WebSocketRequestError} error */
    rejectPending(error) {
        this.pendingRequests.forEach(({ reject, timeout }) => { window.clearTimeout(timeout); reject(error); });
        this.pendingRequests.clear();
    }

    scheduleReconnect() {
        const delay = Math.min(1000 * (2 ** this.reconnectAttempt), 30000);
        this.reconnectAttempt += 1;
        this.reconnectTimer = window.setTimeout(() => this.connect(), delay);
    }
}
