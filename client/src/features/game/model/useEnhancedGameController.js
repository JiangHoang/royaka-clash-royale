import { useEffect } from "react";
import { hasStoredSession } from "@/shared/lib";

/** @param {{ api: { game: ReturnType<import('../api/createGameApi.js').createGameApi> }, isAuthenticated: boolean, navigate: (path: string) => void, handlersRef: { current: Record<string, Function> }, hasLeftGameRef: { current: boolean }, notify: (message: string) => void }} options */
export function useEnhancedGameController({ api, isAuthenticated, navigate, handlersRef, hasLeftGameRef, notify }) {
    useEffect(() => {
        if (!hasStoredSession()) { notify("Session expired. Redirecting to login..."); navigate("/auth"); return; }
        const roomId = localStorage.getItem("room_id");
        if (!roomId) { notify("Room not found. Redirecting to lobby..."); navigate("/lobby"); return; }
        if (!isAuthenticated) return;
        const off = [
            api.game.onTroopChanged((result) => handlersRef.current.troop(result)),
            api.game.onManaChanged((result) => handlersRef.current.mana(result)),
            api.game.onStateChanged((result) => handlersRef.current.state(result)),
            api.game.onGameOver((result, response) => handlersRef.current.gameOver(result, response)),
        ];
        api.game.get(roomId).then((snapshot) => handlersRef.current.snapshot(snapshot)).catch((error) => notify(error.message));
        return () => {
            off.forEach((unsubscribe) => unsubscribe());
            if (!hasLeftGameRef.current) { hasLeftGameRef.current = true; api.game.leave(roomId).catch(() => {}); localStorage.removeItem("room_id"); }
        };
    }, [api, isAuthenticated, navigate, handlersRef, hasLeftGameRef, notify]);
}
