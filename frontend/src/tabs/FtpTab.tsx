import { useState } from 'react';
import type { GeneralTabInformation } from './TabFactory';
import { Alert, AlertDescription, AlertTitle } from '../components/ui/alert';
import { FolderTree, PlusCircle, RefreshCwIcon, ExternalLink } from 'lucide-react';
import { Accordion } from '../components/ui/accordion';
import TabAccordion from '../components/tab-accordion';
import ServerInformation from '../components/server-information';
import { Button } from '../components/ui/button';
import { Dropzone, DropzoneContent, DropzoneEmptyState } from '../components/ui/kibo-ui/dropzone';
import FileTreeView from '../components/filetree/FileTreeView';
import Progress from '@/components/progress';

type FtpTabProps = GeneralTabInformation & {

}

const FtpTab = (props: FtpTabProps) => {
  const [error, setError] = useState<string | null>(null);
  const [droppedFiles, setDroppedFiles] = useState<File[] | undefined>();
  const [uploading, setUploading] = useState<boolean>(false);
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const [uploadedUrl, setUploadedUrl] = useState<string | null>(null);
  const [refreshHandle, setRefreshHandle] = useState<number>(0);

  function handleDropFiles(accFiles: File[]) {
    setDroppedFiles(accFiles);
    setUploadedUrl(null);
    setError(null);
  }

  async function submitUploadFiles() {
    if (!droppedFiles || droppedFiles.length === 0) {
      setError('No files to upload.');
      return;
    }
    setError(null);
    setUploading(true);
    setUploadProgress(0);
    try {
      const file = droppedFiles[0];
      const { uploadFile } = await import('../lib/api');
      const res = await uploadFile(props.id, file, 'ftp', (pct: number) => setUploadProgress(pct));
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
      {error && (
        <Alert className="bg-destructive/10 dark:bg-destructive/20 border-destructive/50 dark:border-destructive/70">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Accordion type="multiple" className="w-full mx-2 space-y-4" defaultValue={['folder_tree']}>
        <ServerInformation id={props.id} reloadTabs={props.reloadTabs} additionalControls={<Button variant="link" asChild><a href={`ftp://localhost`} target="_blank">Open FTP <ExternalLink /></a></Button>} />
        <TabAccordion id="folder_tree" icon={<FolderTree />} title="Folder Tree" tabActions={<Button className="h-8 cursor-pointer" onClick={handleRefresh} title="Refresh" variant={'ghost'}><RefreshCwIcon className="h-4 w-4" /></Button>}>
            <div className="w-full">
            <FileTreeView key={refreshHandle} serverId={props.id} serverType={'ftp'} baseUrl={`ftp://localhost`} />
          </div>
        </TabAccordion>

        <TabAccordion id="upload_resource" icon={<PlusCircle />} title="Upload Resource">
          <Dropzone accept={{ '*/*': [] }} maxFiles={1} maxSize={10 * 1024 * 1024} onDrop={handleDropFiles} onError={console.error} src={droppedFiles}>
            <DropzoneEmptyState />
            <DropzoneContent />
          </Dropzone>
          <Progress active={uploading} value={uploadProgress} className="w-full mb-2 h-2" />
          <Button className="w-full mt-4" disabled={!droppedFiles || droppedFiles.length === 0 || uploading} onClick={submitUploadFiles}>
            {uploading ? 'Uploading...' : 'Upload'}
          </Button>
        </TabAccordion>
      </Accordion>
    </div>
  );
};

export default FtpTab;
