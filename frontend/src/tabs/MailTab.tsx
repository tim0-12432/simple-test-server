import { useState } from "react";
import type { GeneralTabInformation } from "./TabFactory";
import ServerInformation from "../components/server-information";
import { Accordion } from "../components/ui/accordion";
import TabAccordion from "@/components/tab-accordion";
import { MessagesSquare, RefreshCwIcon, FileText } from "lucide-react";
import MailBrowser from "@/components/mail-browser";
import { Button } from "@/components/ui/button";
import { LogsPanel } from "./LogsPanel";

type MailTabProps = GeneralTabInformation & {};

export function MailTab(props: MailTabProps) {
    const [refreshHandle, setRefreshHandle] = useState<number>(0);

    const handleRefresh = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.stopPropagation();
        setRefreshHandle((h) => h + 1);
    };

    const { id, reloadTabs } = props;
    return (
        <div className="w-full h-full flex flex-col items-center gap-4">
            <Accordion type="multiple" className="w-full mx-2 space-y-4" defaultValue={["mail_messages"]}>
                <ServerInformation id={id} reloadTabs={reloadTabs} />
                <TabAccordion id='mail_messages'
                              icon={<MessagesSquare />}
                              title="E-Mails"
                              tabActions={<Button className="h-8 cursor-pointer" onClick={handleRefresh} title="Refresh" variant={'ghost'}><RefreshCwIcon className="h-4 w-4" /></Button>}>
                    <MailBrowser key={refreshHandle} id={id} />
                </TabAccordion>
                <TabAccordion id='server_logs'
                              icon={<FileText />}
                              title="Server Logs"
                              tabActions={<Button className="h-8 cursor-pointer" onClick={handleRefresh} title="Refresh" variant={'ghost'}><RefreshCwIcon className="h-4 w-4" /></Button>}>
                    <div className="w-full">
                        <LogsPanel serverId={id} refreshSignal={refreshHandle} serverType="mail" />
                    </div>
                </TabAccordion>
            </Accordion>
        </div>
    );
}
