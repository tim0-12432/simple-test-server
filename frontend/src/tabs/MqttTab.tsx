import { useEffect, useState } from "react";
import type { GeneralTabInformation } from "./TabFactory";
import { websocketConnect } from "@/lib/api";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { OctagonAlertIcon, FolderTree, ScrollText } from "lucide-react";
import type MqttData from "@/types/MqttData";
import TopicTree from "@/components/topic-tree";
import { Accordion } from "@/components/ui/accordion";
import TabAccordion from "@/components/tab-accordion";
import ServerInformation from "@/components/server-information";

type MqttTabProps = GeneralTabInformation & {

}

const MqttTab = (props: MqttTabProps) => {
    const [error, setError] = useState<string | null>(null);
    const [messages, setMessages] = useState<MqttData[]>([]);

    useEffect(() => {
        if (props.id) {
            websocketConnect(`/protocols/mqtt/${props.id}/messages`, messageHandler, errorHandler);
        }
        setError(null);
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
                <ServerInformation id={props.id} />
                <TabAccordion id='topic_tree'
                              icon={<FolderTree />}
                              title="Topic Tree">
                    <TopicTree messages={messages} />
                </TabAccordion>
                <TabAccordion id='messages'
                              icon={<ScrollText />}
                              title="Message Log">
                </TabAccordion>
            </Accordion>
        </div>
    );
}

export default MqttTab;
