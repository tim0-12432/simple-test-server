import { useEffect } from "react";
import type { GeneralTabInformation } from "./TabFactory";
import { websocketConnect } from "@/lib/api";

type MqttTabProps = GeneralTabInformation & {

}

const MqttTab = (props: MqttTabProps) => {

    useEffect(() => {
        websocketConnect('/protocols/mqtt/e161cbc6c25b8585a2ccffb2d2d2fc0709291a9893c78c247ba318761a8b8374/messages', messageHandler, errorHandler);
    }, []);

    function messageHandler(msg: string) {
        console.log("Received message:", msg);
    }

    function errorHandler(err: Event) {
        console.error("WebSocket error:", err);
    }

    return (
        <>
            MQTT Tab Content
        </>
    );
}

export default MqttTab;
