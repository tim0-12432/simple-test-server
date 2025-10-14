import type { ServerType } from "./Server";


export type Container = {
    container_id: string;
    name: string;
    image: string;
    created_at: number;
    status: number;
    ports: {
        [key: number]: number;
    };
    env: {
        [key: string]: string;
    };
    networks: string[];
    type: ServerType;
}

export { Container };

