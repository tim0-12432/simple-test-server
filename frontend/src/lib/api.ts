
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL ?? (import.meta.env.DEV ? 'http://localhost:8000' : window.location.origin);

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

export type UploadResponse = {
    url: string;
}

export function uploadFile(serverId: string, file: File, serverType: string, onProgress?: (pct: number) => void): Promise<UploadResponse> {
    const url = `${API_URL}/protocols/${serverType}/${encodeURIComponent(serverId)}/upload`;
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', url);
        xhr.onreadystatechange = () => {
            if (xhr.readyState === 4) {
                if (xhr.status >= 200 && xhr.status < 300) {
                    try {
                        const data = JSON.parse(xhr.responseText);
                        resolve(data as UploadResponse);
                    } catch (e) {
                        reject(e);
                    }
                } else {
                    reject(new Error(`HTTP ${xhr.status}: ${xhr.responseText}`));
                }
            }
        };
        xhr.onerror = () => reject(new Error('Network error'));
        if (xhr.upload && typeof onProgress === 'function') {
            xhr.upload.onprogress = (ev) => {
                if (ev.lengthComputable) {
                    const pct = Math.round((ev.loaded / ev.total) * 100);
                    try {
                        onProgress(pct);
                    } catch (e) {
                        // ignore
                    }
                }
            };
        }
        const fd = new FormData();
        fd.append('file', file, file.name);
        xhr.send(fd);
    });
}

import type { FileTreeResponse } from '../types/FileTree';

export async function fetchWebFileTree(serverId: string, path: string | null = null): Promise<FileTreeResponse> {
    const p = path ? `?path=${encodeURIComponent(path)}` : '';
    const url = `/protocols/web/${encodeURIComponent(serverId)}/filetree${p}`;
    const res = await request<FileTreeResponse>("GET", url, undefined);
    return res;
}

// Logs
import type { LogResponse } from '../types/Log';

export async function fetchWebLogs(serverId: string, tail: number = 500): Promise<LogResponse> {
    const url = `/protocols/web/${encodeURIComponent(serverId)}/logs?tail=${tail}`;
    const res = await request<LogResponse>("GET", url, undefined);
    return res;
}

export default request;

