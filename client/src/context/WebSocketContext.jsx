/* eslint-disable react-refresh/only-export-components */
import React, { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from "react";
import { clearStoredSession } from "../utils/session";

const WebSocketContext = createContext(null);

export function WebSocketProvider({ children }) {
    const socketRef = useRef(null);
    const listeners = useRef(new Set());
    const pendingMessages = useRef([]);
    const reconnectTimeout = useRef(null);
    const disposed = useRef(false);
    const [isConnected, setIsConnected] = useState(false);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const wsURL = import.meta.env.VITE_WS_URL || "ws://localhost:8080/ws";

    const sendMessage = useCallback((message) => {
        if (socketRef.current?.readyState !== WebSocket.OPEN) {
            pendingMessages.current.push(message);
            return false;
        }
        socketRef.current.send(JSON.stringify(message));
        return true;
    }, []);

    const connect = useCallback(() => {
        const socket = new WebSocket(wsURL);
        socketRef.current = socket;
        socket.onopen = () => {
            setIsConnected(true);
            const sessionID = localStorage.getItem("session_id");
            if (sessionID) socket.send(JSON.stringify({ type: "authenticate", data: { session_id: sessionID } }));
            pendingMessages.current.splice(0).forEach((message) => socket.send(JSON.stringify(message)));
        };
        socket.onmessage = (event) => {
            let message;
            try { message = JSON.parse(event.data); } catch { return; }
            if (message.type === "login_response" && message.success) {
                localStorage.setItem("session_id", message.data.session_id);
                setIsAuthenticated(true);
            } else if (message.type === "authenticate_response") {
                setIsAuthenticated(Boolean(message.success));
                if (!message.success) clearStoredSession();
            } else if (message.type === "logout_response" && message.success) {
                clearStoredSession();
                setIsAuthenticated(false);
            }
            listeners.current.forEach((listener) => listener(message));
        };
        socket.onclose = () => {
            setIsConnected(false); setIsAuthenticated(false);
            if (!disposed.current) reconnectTimeout.current = window.setTimeout(connect, 3000);
        };
        socket.onerror = (error) => console.error("[WS] Error", error);
    }, [wsURL]);

    useEffect(() => {
        disposed.current = false; connect();
        return () => { disposed.current = true; window.clearTimeout(reconnectTimeout.current); socketRef.current?.close(); };
    }, [connect]);

    const subscribe = useCallback((listener) => { listeners.current.add(listener); return () => listeners.current.delete(listener); }, []);
    const logout = useCallback(() => {
        const sessionID = localStorage.getItem("session_id");
        if (sessionID) sendMessage({ type: "logout", data: { session_id: sessionID } });
        clearStoredSession(); setIsAuthenticated(false);
    }, [sendMessage]);

    const value = useMemo(() => ({ sendMessage, subscribe, logout, isConnected, isAuthenticated }), [sendMessage, subscribe, logout, isConnected, isAuthenticated]);
    return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}

export const useWebSocketContext = () => useContext(WebSocketContext);
