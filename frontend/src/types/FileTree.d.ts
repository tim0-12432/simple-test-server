export type FileTreeEntry = {
    name: string;
    path: string; // relative to webroot
    type: 'file' | 'dir' | 'symlink' | 'unknown';
    size: number;
    modifiedAt: string; // RFC3339
}

export type FileTreeResponse = {
    entries: FileTreeEntry[];
    truncated: boolean;
}
