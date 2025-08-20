
const BACKEND_URL = import.meta.env.DEV ? 'http://localhost:8000' : window.location.origin;

export const API_URL = `${BACKEND_URL}/api/v1`;

export function request<T>(method: string, url: string, body: object|undefined = undefined): Promise<T> {
    const options: RequestInit = {
        method: method.toUpperCase(),
        headers: {
            'Content-Type': 'application/json',
        },
        body: body ? JSON.stringify(body) : null,
    };
    return fetch(API_URL + url, options)
        .then((response) => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json() as Promise<T>;
        })
        .catch((error) => {
            console.error('Fetch error:', error);
            throw error;
        });
}

export function websocketConnect<T>(url: string, onMessage: (msg: T) => void, onError: (err: Event) => void): void {
    const backendUrl = BACKEND_URL.replace(/^http/, 'ws');
    const ws = new WebSocket(`${backendUrl}/api/v1${url}`);
    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            onMessage(data);
        } catch (error) {
            console.error('WebSocket message parsing error:', error);
        }
    };
    ws.onerror = (event) => {
        console.error('WebSocket error:', event);
        onError(event);
    };
}

export default request;
