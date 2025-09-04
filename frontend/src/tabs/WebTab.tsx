import { useState } from "react";
import type { GeneralTabInformation } from "./TabFactory";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { OctagonAlertIcon, FolderTree, PlusCircle, ExternalLink } from "lucide-react";
import { Accordion } from "@/components/ui/accordion";
import TabAccordion from "@/components/tab-accordion";
import ServerInformation from "@/components/server-information";
import { Button } from "@/components/ui/button";
import { Dropzone, DropzoneContent, DropzoneEmptyState } from "@/components/ui/kibo-ui/dropzone";

type WebTabProps = GeneralTabInformation & {

}

const WebTab = (props: WebTabProps) => {
    const [error, setError] = useState<string | null>(null);
    const [droppedFiles, setDroppedFiles] = useState<File[]|undefined>();
    const [port, setPort] = useState<number>(80);

    // TODO: useeffect load filestructure from container
    // TODO: load port for direct link from container info

    function handleDropFiles(accFiles: File[]) {
        console.log("Accepted files:", accFiles);
        setDroppedFiles(accFiles);
    }

    function submitUploadFiles() {
        if (!droppedFiles || droppedFiles.length === 0) {
            setError("No files to upload.");
            return;
        }
        // TODO: POST upload files to container
        setDroppedFiles(undefined);
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
                       defaultValue={['folder_tree']}>
                <ServerInformation id={props.id}
                                   reloadTabs={props.reloadTabs}
                                   additionalControls={<Button variant="link" asChild><a href={`http://localhost:${port}`} target="_blank">Open Webpage <ExternalLink /></a></Button>} />
                <TabAccordion id='folder_tree'
                              icon={<FolderTree />}
                              title="Folder Tree">
                </TabAccordion>
                <TabAccordion id='upload_resource'
                              icon={<PlusCircle />}
                              title="Upload Resource">
                    <Dropzone accept={{ '*/*': [] }}
                              maxFiles={10}
                              maxSize={10 * 1024 * 1024} // 10 MB
                              onDrop={handleDropFiles}
                              onError={console.error}
                              src={droppedFiles}>
                        <DropzoneEmptyState />
                        <DropzoneContent />
                    </Dropzone>
                    <Button className="w-full mt-4" disabled={!droppedFiles || droppedFiles.length === 0} onClick={submitUploadFiles}>Upload</Button>
                </TabAccordion>
            </Accordion>
        </div>
    );
}

export default WebTab;
