import type MqttData from "@/types/MqttData";
import { TreeNode, TreeExpander, TreeIcon, TreeLabel, TreeNodeContent, TreeNodeTrigger, TreeProvider, TreeView } from "./ui/kibo-ui/tree";
import {Folder, FileJson} from "lucide-react"


type TopicTreeProps = {
    messages: MqttData[];
};

type cNode = {
    id: string;
    name: string;
    children?: cNode[];
}

const TopicTree = (props: TopicTreeProps) => {

    function buildTree(messages: MqttData[]): cNode[] {
        const roots: cNode[] = [];
        for (const msg of messages) {
            const parts = msg.topic.split('/');
            parts.push(msg.payload);
            let currentLevel = roots;

            for (const part of parts) {
                let node = currentLevel.find(n => n.name === part);
                if (!node) {
                    node = { id: `${part}-${Math.random()}`, name: part, children: [] };
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

    function renderChildren(children: cNode[], level = 0): React.ReactNode {
        return children.map((child, index) => {
            const childHasChildren = hasChildren(child);
            const isLast = index === children.length - 1;
            return (
                <TreeNode key={child.id} nodeId={child.id} level={level} isLast={isLast} parentPath={[]}>
                    <TreeNodeTrigger>
                        <TreeExpander hasChildren={childHasChildren} />
                        <TreeIcon icon={getIconForNode(child)} hasChildren={childHasChildren} />
                        <TreeLabel>{child.name}</TreeLabel>
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

    const tree = buildTree(props.messages);
    const defaultExpandedIds: string[] = [];
    function collect(nodes: cNode[]) {
        for (const n of nodes) {
            if (hasChildren(n)) {
                defaultExpandedIds.push(n.id);
                if (n.children) collect(n.children);
            }
        }
    }
    collect(tree);

    return (
        <TreeProvider defaultExpandedIds={defaultExpandedIds}>
            <TreeView>
                {renderChildren(tree)}
            </TreeView>
        </TreeProvider>
    )
};

export default TopicTree;
