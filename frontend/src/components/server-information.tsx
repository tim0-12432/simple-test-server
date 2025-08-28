import { Container as Icon } from "lucide-react";
import TabAccordion from "./tab-accordion";
import request from "@/lib/api";
import type { Container } from "@/types/Container";
import { useEffect, useState } from "react";
import { Spinner } from "./ui/kibo-ui/spinner";
import { Button } from "./ui/button";


type ServerInformationProps = {
    id: string;
    reloadTabs: () => void;
}

const ServerInformation = (props: ServerInformationProps) => {
    const [info, setInfo] = useState<Container | null>(null);
    const [stopping, setStopping] = useState<boolean>(false);

    useEffect(() => {
        if (props.id) {
            (async () => {
                try {
                    const container: Container = await request("GET", `/containers/${props.id}`);
                    if (container) {
                        setInfo({...container});
                    } else {
                        console.error("Failed to load container.");
                    }
                } catch (err) {
                    console.error(`Error loading container: ${err instanceof Error ? err.message : "Unknown error"}`);
                }
            })();
        } else {
            setInfo(null);
        }
    }, [props.id]);

    function stopServer() {
        setStopping(true);
        (async () => {
            try {
                await request("DELETE", `/containers/${props.id}`);
            } catch (err) {
                console.error(`Error stopping container: ${err instanceof Error ? err.message : "Unknown error"}`);
            } finally {
                props.reloadTabs();
            }
        })();
    }

    return (
        <TabAccordion id='container_info'
                      icon={<Icon />}
                      title='Container Information'>
            <>
                {
                    info && (
                        <div className="flex">
                            <div className="flex-grow">
                                <label htmlFor="type" className="text-sm font-normal text-muted-foreground">Type:</label>
                                <p id="type" className="font-semibold mb-2">{info.type}</p>
                                <label htmlFor="id" className="text-sm font-normal text-muted-foreground">ID:</label>
                                <p id="id" className="font-semibold mb-2">{info.container_id}</p>
                                <label htmlFor="creation" className="text-sm font-normal text-muted-foreground">Created at:</label>
                                <p id="creation" className="font-semibold mb-2">{new Date(info.created_at).toLocaleString()}</p>
                                <label htmlFor="image" className="text-sm font-normal text-muted-foreground">Docker image:</label>
                                <p id="image" className="font-semibold">{info.image}</p>
                            </div>
                            <div>
                                <Button variant={"outline"} className="cursor-pointer" onClick={stopServer} disabled={stopping}>{stopping ? <>Stopping... <Spinner variant="circle" /></> : "Stop server"}</Button>
                            </div>
                        </div>
                    )
                }
                {
                    !info && (
                        <div>
                            <Spinner variant="circle" />
                        </div>
                    )
                }
            </>
        </TabAccordion>
    );
};

export default ServerInformation;
