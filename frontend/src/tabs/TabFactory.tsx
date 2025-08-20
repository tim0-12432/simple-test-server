import { Alert, AlertDescription } from "@/components/ui/alert";
import type TabType from "@/types/Tab";
import MqttTab from "./MqttTab";
import CreateNewTab from "./CreateNewTab";

export type GeneralTabInformation = {
    id: string;
    type: TabType;
}

const TabFactory = (type: TabType, params: GeneralTabInformation) => {
    let component = null;
    switch (type) {
        case 'MQTT':
            component = <MqttTab {...params} />;
            break;
        case 'create_new':
            component = <CreateNewTab />;
            break;
        default:
            component = <Alert variant='destructive'><AlertDescription>Not yet implemented!</AlertDescription></Alert>
            break;
    }
    return component;
}

export default TabFactory;
