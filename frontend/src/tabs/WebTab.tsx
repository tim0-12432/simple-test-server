import { useState } from "react";
import type { GeneralTabInformation } from "./TabFactory";
import { Alert, AlertDescription, AlertTitle } from "../components/ui/alert";
import { OctagonAlertIcon, FolderTree, PlusCircle, ExternalLink, RefreshCwIcon } from "lucide-react";
import { Accordion } from "../components/ui/accordion";
import TabAccordion from "../components/tab-accordion";
import ServerInformation from "../components/server-information";
import { Button } from "../components/ui/button";
import { Dropzone, DropzoneContent, DropzoneEmptyState } from "../components/ui/kibo-ui/dropzone";
import FileTreeView from "../components/filetree/FileTreeView";
import Progress from "@/components/progress";
import LogsPanel from "./LogsPanel";

type WebTabProps = GeneralTabInformation & {

}

const WebTab = (props: WebTabProps) => {
    const [error, setError] = useState<string | null>(null);
    const [droppedFiles, setDroppedFiles] = useState<File[]|undefined>();
    const [port, _setPort] = useState<number>(80);
    const [uploading, setUploading] = useState<boolean>(false);
    const [uploadProgress, setUploadProgress] = useState<number>(0);
    const [uploadedUrl, setUploadedUrl] = useState<string | null>(null);
    const [refreshHandle, setRefreshHandle] = useState<number>(0);

    // TODO: useeffect load filestructure from container
    // TODO: load port for direct link from container info

    function handleDropFiles(accFiles: File[]) {
        console.log("Accepted files:", accFiles);
        setDroppedFiles(accFiles);
        setUploadedUrl(null);
        setError(null);
    }

    async function submitUploadFiles() {
        if (!droppedFiles || droppedFiles.length === 0) {
            setError("No files to upload.");
            return;
        }
        setError(null);
        setUploading(true);
        setUploadProgress(0);
        try {
            // currently only upload the first file
            const file = droppedFiles[0];
            const { uploadFile } = await import('../lib/api');
            const res = await uploadFile(props.id, file, 'web', (pct) => setUploadProgress(pct));
            setUploadedUrl(res.url);
            setDroppedFiles(undefined);
        } catch (e: any) {
            console.error(e);
            setError(e?.message ?? 'Upload failed');
        } finally {
            setUploading(false);
        }
    }

    const handleRefresh = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.stopPropagation();
        setRefreshHandle((h) => h + 1);
    };

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
                              title="Folder Tree"
                              tabActions={<Button className="h-8 cursor-pointer" onClick={handleRefresh} title="Refresh" variant={"ghost"}><RefreshCwIcon className="h-4 w-4" /></Button>}>
                        <div className="w-full">
                        <FileTreeView key={refreshHandle} serverId={props.id} serverType={'web'} baseUrl={`http://localhost:${port}`} />
                    </div>
                </TabAccordion>
                <TabAccordion id='upload_resource'
                              icon={<PlusCircle />}
                              title="Upload Resource">
                    {/* Upload accordion unchanged */}
                    <Dropzone accept={{ '*/*': [] }}
                              maxFiles={1}
                              maxSize={10 * 1024 * 1024} // 10 MB
                              onDrop={handleDropFiles}
                              onError={console.error}
                              src={droppedFiles}>
                        <DropzoneEmptyState />
                        <DropzoneContent />
                    </Dropzone>
                    <Progress active={uploading} value={uploadProgress} className="w-full mb-2 h-2" />
                    {
                        uploadedUrl ? (
                            <div className="w-full mt-2">
                                <a className="text-primary underline" 
                                   href={uploadedUrl}
                                   target="_blank"
                                   rel="noreferrer">
                                    Open uploaded resource
                                </a>
                            </div>
                        ) : <></>
                    }
                    <Button className="w-full mt-4" disabled={!droppedFiles || droppedFiles.length === 0 || uploading} onClick={submitUploadFiles}>{uploading ? 'Uploading...' : 'Upload'}</Button>
                </TabAccordion>
                <TabAccordion id='access_logs'
                              icon={<OctagonAlertIcon />}
                              title="Access Logs"
                              tabActions={<Button className="h-8 cursor-pointer" onClick={handleRefresh} title="Refresh" variant={"ghost"}><RefreshCwIcon className="h-4 w-4" /></Button>}>
                    <div className="w-full">
                        <LogsPanel serverId={props.id} refreshSignal={refreshHandle} />
                    </div>
                </TabAccordion>
            </Accordion>
        </div>
    );
}

export default WebTab;
