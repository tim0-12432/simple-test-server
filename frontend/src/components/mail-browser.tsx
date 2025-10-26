import { fetchMailMessages } from "@/lib/api";
import { formatBytes } from "@/lib/utils";
import type { MailData } from "@/types/MailData";
import { useEffect, useState } from "react";
import { Progress } from "./progress";

type MailBrowserProps = {
    id: string;
};

const MailBrowser = (props: MailBrowserProps) => {
    const [loading, setLoading] = useState<boolean>(true);
    const [mailMessages, setMailMessages] = useState<MailData[]>([]);
    const [selectedMessage, setSelectedMessage] = useState<MailData | null>(null);

    useEffect(() => {
        setLoading(true);
        (async () => {
            const data = await fetchMailMessages(props.id);
            setMailMessages(data);
            if (data.length > 0) {
                setSelectedMessage(data[0]);
            } else {
                setSelectedMessage(null);
            }
            setLoading(false);
        })();
    }, [props.id]);

    return (
        <>
            <Progress active={loading} className="w-full mb-2 h-2" />
            <div className="w-full h-full flex">
                <div className="w-1/3 h-full overflow-y-auto border-r-2 pr-4">
                    <ul>
                        {mailMessages.map((msg) => (
                            <li key={msg.id}
                                className={`mt-1 p-2 cursor-pointer rounded-md border-2 ${selectedMessage?.id === msg.id ? 'border-gray-600' : 'border-transparent hover:border-blue-500'}`}
                                onClick={() => setSelectedMessage(msg)}>
                                    <div className="font-bold truncate">{msg.content.headers["Subject"]?.[0] || "(No Subject)"}</div>
                                    <div className="text-sm text-muted-foreground truncate">From: {msg.from.name}@{msg.from.domain}</div>
                                    <div className="text-sm text-muted-foreground truncate">Date: {new Date(msg.created).toLocaleString()}</div>
                            </li>
                        ))}
                    </ul>
                </div>
                <div className="w-2/3 h-full overflow-y-auto p-4">
                    {selectedMessage ? (
                        <div>
                            <h2 className="text-xl font-bold mb-2">{selectedMessage.content.headers["Subject"]?.[0] || "(No Subject)"}</h2>
                            <div className="text-sm text-muted-foreground mb-4">
                                Sent: {new Date(selectedMessage.created).toLocaleString()} <br />
                                From: {selectedMessage.from.name}@{selectedMessage.from.domain} <br />
                                To: {selectedMessage.to.map(recipient => `${recipient.name}@${recipient.domain}`).join(", ")} <br />
                            </div>
                            <pre className="whitespace-pre-wrap">{selectedMessage.content.body}</pre>
                            <div className="text-sm text-muted-foreground mt-4">
                                Size: {formatBytes(selectedMessage.content.size)}
                            </div>
                        </div>
                    ) : (
                        <div className="text-muted-foreground text-center">&lt;- Select a message to view its content.</div>
                    )}
                </div>
            </div>
        </>
    );
};

export default MailBrowser;
