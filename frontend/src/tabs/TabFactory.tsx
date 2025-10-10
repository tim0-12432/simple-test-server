import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import type TabType from "@/types/Tab";
import MqttTab from "./MqttTab";
import CreateNewTab from "./CreateNewTab";
import { TriangleAlertIcon } from "lucide-react";
import WebTab from "./WebTab";
import FtpTab from "./FtpTab";
import SmbTab from "./SmbTab";


export type GeneralTabInformation = {
    id: string;
    type: TabType;
    reloadTabs: () => void;
}

const TabFactory = (type: TabType, params: GeneralTabInformation) => {
    let component = null;
    switch (type) {
        case 'create_new':
            component = <CreateNewTab reloadTabs={params.reloadTabs} />;
            break;
        case 'MQTT':
            component = <MqttTab {...params} />;
            break;
        case 'WEB':
            component = <WebTab {...params} />;
            break;
        case 'FTP':
            component = <FtpTab {...params} />;
            break;
        case 'SMB':
            component = <SmbTab {...params} />;
            break;
        default:
            component = (
                <Alert className="bg-amber-500/10 dark:bg-amber-600/30 border-amber-300 dark:border-amber-600/70">
                    <TriangleAlertIcon className="h-4 w-4 !text-amber-500" />
                    <AlertTitle>Warning</AlertTitle>
                    <AlertDescription>
                        Not yet implemented!
                    </AlertDescription>
                </Alert>
            )
            break;
    }
    return component;
}

export default TabFactory;
