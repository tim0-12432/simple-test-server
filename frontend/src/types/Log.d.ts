export type LogLine = {
    ts: string; // RFC3339 timestamp
    line: string;
}

export type LogResponse = {
    lines: LogLine[];
    truncated: boolean;
    container_running: boolean;
}
