import serverTypes from "@/lib/servers";

export type ServerType = typeof serverTypes[number];

export type ServerInformation = {
    name: string;
    image: string;
    ports: number[];
    env: {
        [key: string]: string;
    };
};

export { ServerType };
