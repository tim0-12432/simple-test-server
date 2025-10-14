import { AccordionItem, AccordionTrigger, AccordionContent } from "@/components/ui/accordion";
import type { ReactElement, ReactNode } from "react";


type TabAccordionProps = {
    id: string;
    title: string;
    icon: ReactElement;
    children?: ReactNode;
    tabActions?: ReactElement | ReactElement[] | null;
}

const TabAccordion = (props: TabAccordionProps) => {
    const { id, title, icon, children, tabActions } = props;
    return (
        <AccordionItem value={id}
            className="w-full px-4 border border-border rounded-lg">
            <AccordionTrigger className="cursor-pointer">
                <div className="flex items-start justify-center gap-3 h-6 w-full">
                    {icon}
                    <span>{title}</span>
                    <div className="flex-grow"></div>
                    {tabActions}
                </div>
            </AccordionTrigger>
            <AccordionContent>{children}</AccordionContent>
        </AccordionItem>
    );
};

export default TabAccordion;
