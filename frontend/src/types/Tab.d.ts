import { tabTypes } from '../lib/tabs';

export type TabType = typeof tabTypes[number];

export type Tab = {
    name: string;
    id: string;
    type: TabType;
}

export { TabType };

