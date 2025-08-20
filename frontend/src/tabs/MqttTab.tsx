import type { GeneralTabInformation } from "./TabFactory";

type MqttTabProps = GeneralTabInformation & {

}

const MqttTab = (props: MqttTabProps) => {
    return (
        <>
            MQTT Tab Content
        </>
    );
}

export default MqttTab;
