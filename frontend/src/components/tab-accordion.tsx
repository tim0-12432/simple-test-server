import { AccordionItem, AccordionTrigger, AccordionContent } from "@/components/ui/accordion";
import type { ReactElement } from "react";


type TabAccordionProps = {
    id: string;
    title: string;
    icon: ReactElement;
    children?: ReactElement | ReactElement[] | null;
}

const TabAccordion = (props: TabAccordionProps) => {
    const { id, title, icon, children } = props;
    return (
        <AccordionItem value={id}
            className="w-full px-4 border border-border rounded-lg">
            <AccordionTrigger className="cursor-pointer">
                <div className="flex items-start gap-3">
                    {icon}
                    <span>{title}</span>
                </div>
            </AccordionTrigger>
            <AccordionContent>{children}</AccordionContent>
        </AccordionItem>
    );
};

export default TabAccordion;
