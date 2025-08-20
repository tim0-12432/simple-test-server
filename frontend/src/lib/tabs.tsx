import serverTypes from "./servers";
import { Database, FolderOpen, Globe, Mail, CirclePlus } from "lucide-react"

export const tabTypes = [...serverTypes, 'create_new'] as const;

export function getIconForTabType(type: typeof tabTypes[number]) {
    switch (type) {
        case 'MQTT':
            return <Database />;
        case 'FTP':
            return <FolderOpen />;
        case 'Web':
            return <Globe />;
        case 'SMB':
            return <FolderOpen />;
        case 'Mail':
            return <Mail />;
        case 'create_new':
            return <CirclePlus />;
    }
}

export default tabTypes;
