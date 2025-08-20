import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import Progress from "@/components/progress";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import serverTypes from "@/lib/servers";
import type ServerType from "@/types/Server";
import { useState, useRef, useEffect } from "react";
import { Textarea } from "@/components/ui/textarea";
import { Spinner } from "@/components/ui/kibo-ui/spinner";
import { getIconForTabType } from "@/lib/tabs";


const CreateNewTab = () => {
    const [serverType, setServerType] = useState<ServerType>(serverTypes[0]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const nameRef = useRef<HTMLInputElement>(null);
    const imageRef = useRef<HTMLInputElement>(null);
    const portsRef = useRef<HTMLTextAreaElement>(null);
    const volumesRef = useRef<HTMLTextAreaElement>(null);
    const envRef = useRef<HTMLTextAreaElement>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);
        // Fetch default values for the selected server type
        setTimeout(() => { setLoading(false); }, 1000); // Simulate a network request
    }, [serverType]);

    function handleTypeChange(type: ServerType) {
        setServerType(type);
    }

    function handleSubmit() {
        setLoading(true);
        setError(null);
        setSuccess(null);
        // 1. Validate inputs
        // 2. Send command to server
        setTimeout(() => { setLoading(false); setSuccess("true"); }, 1000); // Simulate a network request
    }

    return (
        <div className="w-full h-full flex flex-col items-center">
            {
                <Progress active={loading} className="w-full mb-2 h-2" />
            }
            <div className="flex flex-col space-y-4 w-full max-w-xs xl:max-w-md 2xl:max-w-lg">
                <div className="w-full space-y-1.5">
                    <Label htmlFor="server-type">Server Type</Label>
                    <Select onValueChange={handleTypeChange} defaultValue={serverTypes[0]}>
                        <SelectTrigger className="w-full">
                            <SelectValue placeholder="Select server type" />
                        </SelectTrigger>
                        <SelectContent>
                            {
                                serverTypes.map((type) => (
                                    <SelectItem key={type} value={type}>
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
                    <Input id="container-name" type="text" placeholder="Name" disabled={loading} />
                    <p className="text-[0.8rem] text-muted-foreground">
                        This will be the name of the container. It is used to identify the container in the list.
                        There is also a autmatic name generation based on the server type e.g. "simple-test-server-{serverType.toLowerCase()}-0".
                    </p>
                </div>
                <div className="w-full space-y-1.5">
                    <Label htmlFor="container-image">Image</Label>
                    <Input id="container-image" type="text" placeholder="Image" disabled={loading} />
                    <p className="text-[0.8rem] text-muted-foreground">
                        This will be the used image for the container. It is used to pull the image from the registry.
                        There is also a predefined image for each server type.
                    </p>
                </div>
                <div className="w-full space-y-1.5">
                    <Label htmlFor="container-ports">Ports</Label>
                    <Textarea id="container-ports" rows={3} placeholder="80:8080" disabled={loading} />
                    <p className="text-[0.8rem] text-muted-foreground">
                        This will be the ports that are exposed by the container. You can specify multiple ports in the format "hostPort:containerPort".
                        For example, "80:8080" will expose port 8080 of the container on port 80 of the host.
                    </p>
                </div>
                <div className="w-full space-y-1.5">
                    <Label htmlFor="container-env">Environment</Label>
                    <Textarea id="container-env" rows={3} placeholder="ENV=PROD" disabled={loading} />
                    <p className="text-[0.8rem] text-muted-foreground">
                        This will be the environment variables that are set in the container.
                        You can specify multiple environment variables in the format "variable=value".
                        For example, "ENV=PROD" will set the environment variable ENV to PROD in the container.
                    </p>
                </div>
                <Button onClick={handleSubmit} disabled={loading} className="w-full">
                    Create {loading && <Spinner variant="circle" />}
                </Button>
            </div>
        </div>
    );
}

export default CreateNewTab;
