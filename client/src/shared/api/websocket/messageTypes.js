
export const MESSAGE_TYPES = Object.freeze({
    REGISTER: "register",
    LOGIN: "login",
    AUTHENTICATE: "authenticate",
    LOGOUT: "logout",
    GET_USER: "get_user",
    GET_DECK: "get_desk",
    FIND_MATCH: "find_match",
    GET_GAME: "get_game",
    ATTACK: "attack",
    HEAL: "heal",
    SKIP_TURN: "skip_turn",
    PLAY_AGAIN: "play_again",
    LEAVE_GAME: "leave_game",
    SELECT_TROOP: "select_troop",
    ERROR: "error",
    REGISTER_RESPONSE: "register_response",
    LOGIN_RESPONSE: "login_response",
    AUTHENTICATE_RESPONSE: "authenticate_response",
    LOGOUT_RESPONSE: "logout_response",
    USER_RESPONSE: "user_response",
    DECK_RESPONSE: "deck_response",
    FIND_MATCH_RESPONSE: "find_match_response",
    MATCH_FOUND: "match_found",
    MATCH_TIMEOUT: "match_timeout",
    GAME_RESPONSE: "game_response",
    ATTACK_RESPONSE: "attack_response",
    HEAL_RESPONSE: "heal_response",
    SKIP_TURN_RESPONSE: "skip_turn_response",
    PLAY_AGAIN_RESPONSE: "play_again_response",
    LEAVE_GAME_RESPONSE: "leave_game_response",
    TROOP_RESPONSE: "troop_response",
    MANA_UPDATE: "mana_update",
    GAME_STATE: "game_state",
    GAME_OVER_RESPONSE: "game_over_response",
});

/** @typedef {typeof MESSAGE_TYPES[keyof typeof MESSAGE_TYPES]} MessageType */

const messageTypeSet = new Set(Object.values(MESSAGE_TYPES));

/** @param {unknown} value @returns {value is MessageType} */
export const isMessageType = (value) => typeof value === "string" && messageTypeSet.has(/** @type {MessageType} */ (value));
