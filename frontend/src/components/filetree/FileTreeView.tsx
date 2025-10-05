'use client';

import { useEffect } from 'react';
import {
  TreeProvider,
  TreeView,
  TreeNode,
  TreeNodeTrigger,
  TreeNodeContent,
  TreeExpander,
  TreeIcon,
  TreeLabel,
} from '../ui/kibo-ui/tree';
import { useFileTree } from '@/lib/useFileTree';
import { ExternalLink, FileIcon, FolderIcon } from 'lucide-react';

type FileTreeViewProps = {
  serverId: string;
  baseUrl?: string; // e.g. http://localhost:8080
};

function humanSize(n: number) {
  if (n < 1024) return `${n} B`;
  if (n < 1024 * 1024) return `${Math.round((n / 1024) * 10) / 10} KB`;
  return `${Math.round((n / (1024 * 1024)) * 10) / 10} MB`;
}

function NodeChildren({ path, serverId, baseUrl }: { path: string | null; serverId: string; baseUrl?: string }) {
  const { getChildren, getCached, loadingPaths } = useFileTree(serverId);

  const key = path || '';
  const cached = getCached(path);
  const isLoading = !!loadingPaths[key];

  useEffect(() => {
    let mounted = true;
    (async () => {
      if (!cached && mounted) {
        try {
          await getChildren(path ?? null);
        } catch (e) {
          // ignore - parent will show errors elsewhere
        }
      }
    })();
    return () => {
      mounted = false;
    };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [path]);

  if (isLoading && !cached) {
    return <div className="px-4 py-2 text-sm text-muted-foreground">Loading...</div>;
  }

  const entries = cached?.entries ?? [];

  if (entries.length === 0) {
    return <div className="px-4 py-2 text-sm text-muted-foreground">Empty folder</div>;
  }

  return (
    <div>
      {entries.map((e, idx) => {
        const nodeId = e.path || e.name;
        const hasChildren = e.type === 'dir';
        const childPath = e.path;
        return (
          <TreeNode key={nodeId} nodeId={nodeId} level={1} isLast={idx === entries.length - 1} parentPath={[]}>
            <div className="flex items-center gap-2">
              <TreeExpander hasChildren={hasChildren} />
              <TreeIcon hasChildren={hasChildren} icon={hasChildren ? <FolderIcon /> : <FileIcon />} />
              <TreeNodeTrigger>
                <div className="flex w-full items-center gap-2">
                  <TreeLabel>
                    {e.name}
                  </TreeLabel>
                  <span className="text-xs text-muted-foreground">{e.type === 'file' ? humanSize(e.size) : ''}</span>
                </div>
              </TreeNodeTrigger>
            </div>

            <TreeNodeContent hasChildren={true} className="pl-6">
              {hasChildren ? (
                <NodeChildren path={childPath} serverId={serverId} baseUrl={baseUrl} />
              ) : (
                <>
                {baseUrl ? (
                  <div className="px-8 py-1">
                    <TreeNode level={2}>
                      <div className="flex items-center gap-2">
                        <TreeIcon icon={<ExternalLink className="h-4 w-4" />} />
                        <TreeLabel className="group relative mx-1 px-3 py-2">
                          <a className="text-primary underline px-3 py-2 rounded-md transition-all duration-200 hover:bg-accent/50 bg-accent/0" href={`${baseUrl}/${e.path}`} target="_blank" rel="noreferrer">
                            Open
                          </a>
                        </TreeLabel>
                      </div>
                    </TreeNode>
                  </div>
                ) : null}
                </>
              )}
            </TreeNodeContent>
          </TreeNode>
        );
      })}
    </div>
  );
}

export default function FileTreeView({ serverId, baseUrl }: FileTreeViewProps) {
  const { getChildren, getCached } = useFileTree(serverId);

  useEffect(() => {
    // load root on mount
    (async () => {
      try {
        await getChildren(null);
      } catch (e) {
        // ignore here
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [serverId]);

  const rootCached = getCached(null);

  return (
    <TreeProvider defaultExpandedIds={[]} showIcons selectable indent={18} animateExpand>
      <TreeView>
        {rootCached ? (
          <NodeChildren path={null} serverId={serverId} baseUrl={baseUrl} />
        ) : (
          <div className="px-4 py-2 text-sm text-muted-foreground">Loading...</div>
        )}
      </TreeView>
    </TreeProvider>
  );
}
