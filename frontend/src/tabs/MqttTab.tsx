import { useEffect, useRef, useState } from "react";
import type { GeneralTabInformation } from "./TabFactory";
import { websocketConnect } from "@/lib/api";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { OctagonAlertIcon, FolderTree, ScrollText } from "lucide-react";
import type MqttData from "@/types/MqttData";
import TopicTree from "@/components/topic-tree";
import MessageLog from "@/components/message-log/MessageLog";
import { Accordion } from "@/components/ui/accordion";
import TabAccordion from "@/components/tab-accordion";
import ServerInformation from "@/components/server-information";

type MqttTabProps = GeneralTabInformation & {

}

const MqttTab = (props: MqttTabProps) => {
    const [error, setError] = useState<string | null>(null);
    const [messages, setMessages] = useState<MqttData[]>([]);

    const wsRef = useRef<WebSocket | null>(null);
    const [connected, setConnected] = useState(false);

    useEffect(() => {
        setError(null);
        // close previous socket if any
        if (wsRef.current) {
            try { wsRef.current.close(); } catch (e) { /* ignore */ }
            wsRef.current = null;
            setConnected(false);
        }

        if (props.id) {
            const ws = websocketConnect(`/protocols/mqtt/${props.id}/messages`, messageHandler, errorHandler);
            wsRef.current = ws;
            ws.onopen = () => setConnected(true);
            ws.onclose = () => setConnected(false);
            ws.onerror = () => setConnected(false);
        }

        return () => {
            if (wsRef.current) {
                try { wsRef.current.close(); } catch (e) { /* ignore */ }
                wsRef.current = null;
            }
        };
    }, [props.id]);

    function messageHandler(msg: MqttData) {
        setMessages(prevMessages => [...prevMessages, msg]);
    }

    function errorHandler(err: Event) {
        setError(`WebSocket error: ${err instanceof Error ? err.message : "Unknown error"}`);
    }

    return (
        <div className="w-full h-full flex flex-col items-center gap-4">
            {
                error && (
                    <Alert className="bg-destructive/10 dark:bg-destructive/20 border-destructive/50 dark:border-destructive/70">
                        <OctagonAlertIcon className="h-4 w-4 !text-destructive" />
                        <AlertTitle>Error</AlertTitle>
                        <AlertDescription>
                            {error}
                        </AlertDescription>
                    </Alert>
                )
            }
            <Accordion type="multiple"
                       className="w-full mx-2 space-y-4"
                       defaultValue={['topic_tree']}>
                <ServerInformation id={props.id} reloadTabs={props.reloadTabs} />
                <TabAccordion id='topic_tree'
                              icon={<FolderTree />}
                              title="Topic Tree">
                    <TopicTree messages={messages} />
                </TabAccordion>
                <TabAccordion id='messages'
                              icon={<ScrollText />}
                              title="Message Log">
                    <MessageLog messages={messages} />
                </TabAccordion>
            </Accordion>
        </div>
    );
}

export default MqttTab;
