
export const hasStoredSession = () => Boolean(localStorage.getItem("session_id"));
export const clearStoredSession = () => {
    localStorage.removeItem("session_id");
    localStorage.removeItem("username");
    localStorage.removeItem("room_id");
};
