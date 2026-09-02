/* eslint-disable react-refresh/only-export-components */
import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { createApi } from "./createApi.js";
import { MESSAGE_TYPES, WebSocketTransport } from "@/shared/api/websocket";
import { clearStoredSession } from "@/shared/lib";
import { WS_URL } from "@/shared/config";

/** @typedef {'connecting'|'connected'|'disconnected'|'error'} ConnectionState */
/** @typedef {{ api: ReturnType<import('./createApi.js').createApi>, connectionState: ConnectionState, isConnected: boolean, isAuthenticated: boolean }} WebSocketContextValue */
const WebSocketContext = createContext(/** @type {WebSocketContextValue | null} */ (null));

export function WebSocketProvider({ children }) {
    const [connectionState, setConnectionState] = useState(/** @type {ConnectionState} */ ("disconnected"));
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const transport = useMemo(() => new WebSocketTransport(WS_URL, {
        onStateChange: setConnectionState,
        onAuthenticated: setIsAuthenticated,
    }), []);
    const api = useMemo(() => createApi(transport), [transport]);

    useEffect(() => {
        transport.connect();
        return () => transport.close();
    }, [transport]);

    useEffect(() => transport.on(MESSAGE_TYPES.LOGIN_RESPONSE, (_session, response) => {
        if (response.success) setIsAuthenticated(true);
    }), [transport]);

    useEffect(() => transport.on(MESSAGE_TYPES.LOGOUT_RESPONSE, (_data, response) => {
        if (response.success) {
            clearStoredSession();
            setIsAuthenticated(false);
        }
    }), [transport]);

    const value = useMemo(() => ({
        api,
        connectionState,
        isConnected: connectionState === "connected",
        isAuthenticated,
    }), [api, connectionState, isAuthenticated]);

    return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}

export const useWebSocketContext = () => {
    const context = /** @type {WebSocketContextValue | null} */ (useContext(WebSocketContext));
    if (!context) throw new Error("useWebSocketContext must be used inside WebSocketProvider");
    return context;
};
