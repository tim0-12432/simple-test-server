import type MqttData from "@/types/MqttData";
import { TreeExpander, TreeIcon, TreeLabel, TreeNodeContent, TreeNodeTrigger, TreeProvider, TreeView } from "./ui/kibo-ui/tree";
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

    function renderChildren(children: cNode[]): React.ReactNode {
        return children.map((child) => (
            <>
                <TreeNodeTrigger key={child.id}>
                    <TreeExpander />
                    <TreeIcon icon={getIconForNode(child)} />
                    <TreeLabel>{child.name}</TreeLabel>
                </TreeNodeTrigger>
                {
                    hasChildren(child) ? (
                        <TreeNodeContent>
                            {renderChildren(child.children!)}
                        </TreeNodeContent>
                    ) : null
                }
            </>
        ));
    }

    return (
        <TreeProvider>
            <TreeView>
                {renderChildren(buildTree(props.messages))}
            </TreeView>
        </TreeProvider>
    )
};

export default TopicTree;
