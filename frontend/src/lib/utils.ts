import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function hashCode(s: string) {
  const bytes = new TextEncoder().encode(s);
  return window.crypto.subtle.digest('SHA-1', bytes).then(hashBuffer => {
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  });
}

export function formatBytes(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const dm = decimals < 0 || i === 0 ? 0 : decimals;
  const value = (bytes / Math.pow(k, i)).toFixed(dm);
  return `${value} ${sizes[i]}`;
}
