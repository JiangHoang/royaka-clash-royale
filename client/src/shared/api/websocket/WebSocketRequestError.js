
export class WebSocketRequestError extends Error {
    /**
     * @param {string} message
     * @param {string} [code]
     * @param {import('./contracts.js').WsResponse<unknown> | null} [response]
     */
    constructor(message, code = "request_failed", response = null) {
        super(message);
        this.name = "WebSocketRequestError";
        this.code = code;
        this.response = response;
    }
}
