import type { GeneralTabInformation } from "./TabFactory";
import ServerInformation from "../components/server-information";
import { Accordion } from "../components/ui/accordion";

type MailTabProps = GeneralTabInformation & {};

const MailTab = (props: MailTabProps) => {
    return (
        <div className="w-full h-full flex flex-col items-center gap-4">
            <Accordion type="multiple" className="w-full mx-2 space-y-4" defaultValue={["container_info"]}>
                <ServerInformation id={props.id} reloadTabs={props.reloadTabs} />
            </Accordion>
        </div>
    );
};

export default MailTab;
