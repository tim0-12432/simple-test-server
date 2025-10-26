import { useEffect, useState } from 'react';
import type { GeneralTabInformation } from './TabFactory';
import { useOtelTelemetry } from '@/lib/useOtelTelemetry';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { OctagonAlertIcon, CircleCheck, CircleX, MessageSquareTextIcon, FunnelX, MessageSquareX } from 'lucide-react';
import ServerInformation from '@/components/server-information';
import { Button } from '@/components/ui/button';
import { Accordion } from '@/components/ui/accordion';
import TabAccordion from '@/components/tab-accordion';
import { Input } from '@/components/ui/input';

type OtelTabProps = GeneralTabInformation & {};

export function OtelTab(props: OtelTabProps) {
    const [error, setError] = useState<string | null>(null);
    const { messages, connected, error: wsError, clear } = useOtelTelemetry(props.id);
    const [filteredMessages, setFilteredMessages] = useState<string[]>(messages);
    const [filterText, setFilterText] = useState<string>('');

    // reflect websocket error into local error state
    if (wsError && wsError !== error) {
        setError(wsError);
    }

    useEffect(() => {
        if (filterText.trim() === '') {
            setFilteredMessages(messages);
        } else {
            const lowerFilter = filterText.toLowerCase();
            const filtered = messages.filter(m => m.toLowerCase().includes(lowerFilter));
            setFilteredMessages(filtered);
        }
    }, [messages, filterText]);

    function onFilterChange(e: React.ChangeEvent<HTMLInputElement>) {
        setFilterText(e.target.value);
    }

    const ServerControls = (
        <>
            <div className="flex items-center justify-end text-sm py-2 gap-1">
                {connected ? (
                    <>
                        <CircleCheck className="inline h-4 w-4 mr-1 text-green-500" />
                        Connected
                    </>
                ) : (
                    <>
                        <CircleX className="inline h-4 w-4 mr-1 text-red-500" />
                        Disconnected
                    </>
                )}
            </div>
            <div className="flex items-center justify-end">
                <Button variant="outline" className="flex items-center justify-center gap-2 cursor-pointer"
                        onClick={() => clear()}
                        disabled={messages.length === 0}>
                    Clear Messages
                    <MessageSquareX className="h-4 w-4 mr-1" />
                </Button>
            </div>
        </>
    );

    return (
        <div className="w-full h-full flex flex-col items-center gap-4">
            {error && (
                <Alert className="bg-destructive/10 dark:bg-destructive/20 border-destructive/50 dark:border-destructive/70">
                    <OctagonAlertIcon className="h-4 w-4 !text-destructive" />
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription>{error}</AlertDescription>
                </Alert>
            )}

            <Accordion type="multiple"
                       className="w-full mx-2 space-y-4"
                       defaultValue={['telemetry_messages']}>
                <ServerInformation id={props.id} reloadTabs={props.reloadTabs} additionalControls={ServerControls} />
                <TabAccordion id='telemetry_messages'
                                icon={<MessageSquareTextIcon />}
                                title="Telemetry Messages">
                    <div className="py-2 flex gap-2 items-center justify-center">
                        Filter:
                        <Input type="text"
                               placeholder="Type to filter messages..."
                               value={filterText}
                               onChange={onFilterChange}
                               disabled={messages.length === 0}
                               className="w-full px-3 py-2 border border-border rounded outline-none bg-input text-input-foreground text-sm" />
                        <Button className="cursor-pointer" variant="ghost" disabled={messages.length === 0 || filterText === ''} onClick={() => setFilterText('')}>
                            <FunnelX className="h-4 w-4" />
                        </Button>
                    </div>
                    {messages.length === 0 ? (
                        <div className="text-sm text-muted-foreground">No telemetry received yet</div>
                    ) : (
                        <ul>
                            {filteredMessages.map((m, idx) => {
                                let formatted = m;
                                try {
                                    const parsed = JSON.parse(m);
                                    formatted = JSON.stringify(parsed, null, 2);
                                } catch {
                                    // keep original
                                }
                                return (
                                    <li key={idx} className="py-2 pb-3 border-b last:border-0 border-border">
                                        <pre className="text-sm px-4 py-2 whitespace-pre-wrap">{formatted}</pre>
                                    </li>
                                );
                            })}
                        </ul>
                    )}
                </TabAccordion>
            </Accordion>
        </div>
    );
}
