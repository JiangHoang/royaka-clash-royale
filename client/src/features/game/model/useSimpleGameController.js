import { useEffect } from "react";
import { hasStoredSession } from "@/shared/lib";

/**
 * Owns the WebSocket lifecycle for the simple game screen.
 * @param {{ api: { game: ReturnType<import('../api/createGameApi.js').createGameApi> }, isAuthenticated: boolean, navigate: (path: string) => void, handlersRef: { current: Record<string, Function> }, hasLeftGameRef: { current: boolean }, damageTimeoutRef: { current: number | null }, healTimeoutRef: { current: number | null }, notify: (message: string) => void }} options
 */
export function useSimpleGameController({ api, isAuthenticated, navigate, handlersRef, hasLeftGameRef, damageTimeoutRef, healTimeoutRef, notify }) {
    useEffect(() => {
        if (!hasStoredSession()) { notify("Session expired. Redirecting to login..."); navigate("/auth"); return; }
        const roomId = localStorage.getItem("room_id");
        if (!roomId) { notify("Room not found. Redirecting to lobby..."); navigate("/lobby"); return; }
        if (!isAuthenticated) return;
        const off = [
            api.game.onAttack((result) => handlersRef.current.attack(result)),
            api.game.onHeal((result) => handlersRef.current.heal(result)),
            api.game.onTurnChanged((result) => handlersRef.current.turn(result)),
            api.game.onGameOver((result, response) => handlersRef.current.gameOver(result, response)),
        ];
        api.game.get(roomId).then((snapshot) => handlersRef.current.snapshot(snapshot)).catch((error) => notify(error.message));
        return () => {
            off.forEach((unsubscribe) => unsubscribe());
            if (!hasLeftGameRef.current) { hasLeftGameRef.current = true; api.game.leave(roomId).catch(() => {}); localStorage.removeItem("room_id"); }
            // Popup timers are intentionally read at cleanup time.
            if (damageTimeoutRef.current !== null) clearTimeout(damageTimeoutRef.current);
            if (healTimeoutRef.current !== null) clearTimeout(healTimeoutRef.current);
        };
    }, [api, isAuthenticated, navigate, handlersRef, hasLeftGameRef, damageTimeoutRef, healTimeoutRef, notify]);
}
