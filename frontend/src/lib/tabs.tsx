import serverTypes from "./servers";
import { Database, FolderOpen, Globe, Mail, CirclePlus, Telescope, type LucideProps } from "lucide-react"

export const tabTypes = [...serverTypes, 'create_new'] as const;

export function getIconForTabType(type: typeof tabTypes[number], params?: LucideProps) {
    switch (type) {
        case 'MQTT':
            return <Database {...params} />;
        case 'FTP':
            return <FolderOpen {...params} />;
        case 'WEB':
            return <Globe {...params} />;
        case 'SMB':
            return <FolderOpen {...params} />;
        case 'MAIL':
            return <Mail {...params} />;
        case 'OTEL':
            return <Telescope {...params} />;
        case 'create_new':
            return <CirclePlus {...params} />;
    }

}

export default tabTypes;
