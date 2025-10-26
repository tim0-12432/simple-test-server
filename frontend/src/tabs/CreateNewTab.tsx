import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Progress } from "@/components/progress";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import serverTypes from "@/lib/servers";
import type { ServerType } from "@/types/Server";
import { useState, useRef, useEffect } from "react";
import { Textarea } from "@/components/ui/textarea";
import { Spinner } from "@/components/ui/kibo-ui/spinner";
import { getIconForTabType } from "@/lib/tabs";
import request, { API_URL } from "@/lib/api";
import type { ServerInformation } from "@/types/Server";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { CircleCheckBigIcon, OctagonAlertIcon } from "lucide-react";

type LoadingState = {
    message: string;
    percent: number;
}

type CreateNewTabProps = {
    reloadTabs: () => void;
}

export function CreateNewTab(props: CreateNewTabProps) {
    const [serverType, setServerType] = useState<ServerType>(serverTypes[0]);
    const [loading, setLoading] = useState(false);
    const [loadingState, setLoadingState] = useState<LoadingState | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const nameRef = useRef<HTMLInputElement>(null);
    const imageRef = useRef<HTMLInputElement>(null);
    const portsRef = useRef<HTMLTextAreaElement>(null);
    const envRef = useRef<HTMLTextAreaElement>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);
        setSuccess(null);
        (async () => {
            try {
                const server: ServerInformation = await request("GET", `/servers/${serverType.toUpperCase()}`);
                if (server) {
                    if (nameRef.current) {
                        nameRef.current.value = "";
                    }
                    if (imageRef.current) {
                        imageRef.current.value = server.image;
                    }
                    if (portsRef.current) {
                        portsRef.current.value = server.ports.map(p => `${p}:${p}`).join("\n");
                    }
                    if (envRef.current) {
                        envRef.current.value = Object.entries(server.env).map(([key, value]) => `${key}=${value}`).join("\n");
                    }
                } else {
                    setError("Failed to load server definition.");
                }
            } catch (err) {
                setError(`Error loading server definition: ${err instanceof Error ? err.message : "Unknown error"}`);
            }
            setLoading(false);
        })();
    }, [serverType]);

    function handleTypeChange(type: ServerType) {
        setServerType(type);
    }

    function handleSubmit() {
        setLoading(true);
        setError(null);
        setSuccess(null);
        (async () => {
            function closeEventSource(eventSource: EventSource) {
                eventSource.close();
                if (!error) {
                    setSuccess("Server created successfully.");
                }
                setLoading(false);
                props.reloadTabs();
            }

            const name = nameRef.current?.value;
            const image = imageRef.current?.value;
            const ports = portsRef.current?.value.split("\n").map(p => {
                const parts = p.split(":");
                return parts.length === 2 ? { [parseInt(parts[0]).toString()]: parseInt(parts[1]) } : null;
            }).filter(p => p !== null);
            const env = envRef.current?.value.split("\n").reduce((acc, line) => {
                const parts = line.split("=");
                if (parts.length === 2) {
                    acc[parts[0].trim()] = parts[1].trim();
                }
                return acc;
            }, {} as Record<string, string>);
            const {reqId} = await request("POST", `/servers/${serverType.toUpperCase()}`, {
                name, image, ports, env
            }) as { reqId: string };

            const eventSource = new EventSource(`${API_URL}/servers/progress/${reqId}`);
            eventSource.onmessage = (event) => {
                // {"percent":50,"message":"pull failed: docker pull failed: exit status 1 - Error response from daemon: Head "https://ghcr.io/v2/servercontainers/mailbox/manifests/latest": denied","error":true}
                const data = JSON.parse(event.data) as {percent: number; message: string; error: boolean};
                setLoadingState(data);
                if (data.error) {
                    setError(data.message);
                    setLoading(false);
                }
                if (data.percent >= 100 || data.error) {
                    closeEventSource(eventSource);
                }
            }
            eventSource.onerror = () => closeEventSource(eventSource);
        })();
    }

    return (
        <div className="w-full h-full flex flex-col items-center">
            {
                <Progress active={loading} value={loadingState?.percent} className="w-full mb-2 h-2" />
            }
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
            {
                success && (
                    <Alert className="bg-emerald-500/10 dark:bg-emerald-600/30 border-emerald-300 dark:border-emerald-600/70">
                        <CircleCheckBigIcon className="h-4 w-4 !text-emerald-500" />
                        <AlertTitle>Successful</AlertTitle>
                        <AlertDescription>
                            {success}
                        </AlertDescription>
                    </Alert>
                )
            }
            {
                !error && !success && (
                    <div className="flex flex-col space-y-4 w-full max-w-xs xl:max-w-md 2xl:max-w-lg">
                        <div className="w-full space-y-1.5">
                            <Label htmlFor="server-type">Server Type</Label>
                            <Select disabled={loading} onValueChange={handleTypeChange} defaultValue={serverTypes[0]}>
                                <SelectTrigger className="w-full cursor-pointer">
                                    <SelectValue placeholder="Select server type" />
                                </SelectTrigger>
                                <SelectContent>
                                    {
                                        serverTypes.map((type) => (
                                            <SelectItem key={type} value={type} className="cursor-pointer">
                                                {getIconForTabType(type)} {type}
                                            </SelectItem>
                                        ))
                                    }
                                </SelectContent>
                            </Select>
                        </div>
                        <hr className="border-border" />
                        <div className="w-full space-y-1.5">
                            <Label htmlFor="container-name">Name</Label>
                            <Input id="container-name" type="text" placeholder={`simple-test-server-${serverType.toLowerCase()}-0`} disabled={loading} ref={nameRef} />
                            <p className="text-[0.8rem] text-muted-foreground">
                                This will be the name of the container. It is used to identify the container in the list.
                                There is also a autmatic name generation based on the server type e.g. "simple-test-server-{serverType.toLowerCase()}-0".
                            </p>
                        </div>
                        <div className="w-full space-y-1.5">
                            <Label htmlFor="container-image">Image</Label>
                            <Input id="container-image" type="text" placeholder={`image-${serverType.toLowerCase()}:latest`} disabled={loading} ref={imageRef} />
                            <p className="text-[0.8rem] text-muted-foreground">
                                This will be the used image for the container. It is used to pull the image from the registry.
                                There is also a predefined image for each server type.
                            </p>
                        </div>
                        <div className="w-full space-y-1.5">
                            <Label htmlFor="container-ports">Ports</Label>
                            <Textarea id="container-ports" rows={3} placeholder="80:8080" disabled={loading} ref={portsRef} />
                            <p className="text-[0.8rem] text-muted-foreground">
                                This will be the ports that are exposed by the container. You can specify multiple ports in the format "hostPort:containerPort".
                                For example, "80:8080" will expose port 8080 of the container on port 80 of the host.
                            </p>
                        </div>
                        <div className="w-full space-y-1.5">
                            <Label htmlFor="container-env">Environment</Label>
                            <Textarea id="container-env" rows={3} placeholder="ENV=PROD" disabled={loading} ref={envRef} />
                            <p className="text-[0.8rem] text-muted-foreground">
                                This will be the environment variables that are set in the container.
                                You can specify multiple environment variables in the format "variable=value".
                                For example, "ENV=PROD" will set the environment variable ENV to PROD in the container.
                            </p>
                        </div>
                        <Button onClick={handleSubmit} disabled={loading} className="w-full cursor-pointer">
                            {loading ? <>{loadingState?.message} <Spinner variant="circle" /></> : "Create"}
                        </Button>
                    </div>
                )
            }
        </div>
    );
}
