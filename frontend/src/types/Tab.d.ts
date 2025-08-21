import { tabTypes } from '../lib/tabs';

export type TabType = typeof tabTypes[number];

export default TabType;

export type Tab = {
    name: string;
    id: string;
    type: TabType;
}
