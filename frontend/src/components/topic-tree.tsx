import type MqttData from "@/types/MqttData";
import { TreeNode, TreeExpander, TreeIcon, TreeNodeContent, TreeNodeTrigger, TreeProvider, TreeView, type TreeLabelProps, TreeLabel } from "./ui/kibo-ui/tree";
import {Folder, FileJson} from "lucide-react"
import { useEffect, useState } from "react";
import { hashCode } from "@/lib/utils";


type TopicTreeProps = {
    messages: MqttData[];
};

type cNode = {
    id: string;
    name: string;
    children?: cNode[];
}

const TopicTree = (props: TopicTreeProps) => {
    const [loading, setLoading] = useState(true);
    const [tree, setTree] = useState<cNode[]>([]);
    const [defaultExpandedIds, setDefaultExpandedIds] = useState<string[]>([]);
    const [expanded, setExpanded] = useState<string[]|null>(null);

    async function buildTree(messages: MqttData[]): Promise<cNode[]> {
        const roots: cNode[] = [];
        for (const msg of messages) {
            const parts = msg.topic.split('/');
            parts.push(msg.payload);
            let currentLevel = roots;

            for (const part of parts) {
                const pathUntilNode = parts.slice(0, parts.indexOf(part) + 1).join('/').toLowerCase();
                let node = currentLevel.find(n => n.name === part);
                if (!node) {
                    node = { id: await hashCode(pathUntilNode), name: part, children: [] };
                    currentLevel.push(node);
                }
                currentLevel = node.children!;
            }
        }
        return roots;
    }

    function hasChildren(node: cNode): boolean {
        return node.children && node.children.length > 0 || false;
    }

    function getIconForNode(node: cNode): React.ReactNode {
        if (hasChildren(node)) {
            return <Folder className="h-4 w-4" />;
        }
        return <FileJson className="h-4 w-4" />;
    }

    function onClickNode(nodeId: string) {
        setExpanded(prev => {
            let newExpanded = new Set(prev);
            if (prev == null) {
                newExpanded = new Set(defaultExpandedIds);
            }
            if (newExpanded.has(nodeId)) {
                newExpanded.delete(nodeId);
            } else {
                newExpanded.add(nodeId);
            }
            return Array.from(newExpanded);
        });
    }

    function renderChildren(children: cNode[], level = 0): React.ReactNode {
        return children.map((child, index) => {
            const childHasChildren = hasChildren(child);
            const isLast = index === children.length - 1;
            return (
                <TreeNode key={child.id} nodeId={child.id} level={level} isLast={isLast} parentPath={[]}>
                    <TreeNodeTrigger>
                        <TreeExpander hasChildren={childHasChildren} onClick={() => onClickNode(child.id)} />
                        <TreeIcon icon={getIconForNode(child)} hasChildren={childHasChildren} />
                        {
                            !childHasChildren
                            ? <PayloadTreeLabel>{child.name}</PayloadTreeLabel>
                            : <TreeLabel>{child.name}</TreeLabel>
                        }
                    </TreeNodeTrigger>
                    {childHasChildren ? (
                        <TreeNodeContent hasChildren={true}>
                            {renderChildren(child.children!, level + 1)}
                        </TreeNodeContent>
                    ) : null}
                </TreeNode>
            );
        });
    }

    useEffect(() => {
        setLoading(true);
        (async () => {
            const builtTree = await buildTree(props.messages);
            const defaultIds: string[] = [];
            function collect(nodes: cNode[], level: number) {
                if (level === 3) {
                    return;
                }
                for (const n of nodes) {
                    if (hasChildren(n)) {
                        defaultIds.push(n.id);
                        if (n.children) collect(n.children, level + 1);
                    }
                }
            }
            collect(builtTree, 1);
            setTree(builtTree);
            setDefaultExpandedIds(defaultIds);
            setLoading(false);
        })();
    }, [props.messages]);

    if (!loading) {
        if (tree.length === 0) {
            return <div className="text-sm text-muted-foreground">No topics available.</div>;
        }
        return (
            <TreeProvider defaultExpandedIds={expanded == null ? defaultExpandedIds : expanded}>
                <TreeView>
                    {renderChildren(tree)}
                </TreeView>
            </TreeProvider>
        );
    }
    return <div>Loading...</div>;
};

export default TopicTree;

const PayloadTreeLabel = (props: TreeLabelProps) => (
  <span className='font flex-1 wrap-anywhere text-sm' {...props} />
);
