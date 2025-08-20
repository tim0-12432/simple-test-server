import { tabTypes } from '../lib/tabs';

export type TabType = typeof tabTypes[number];

export default TabType;

type TabTypesForIndex = Exclude<TabType, 'create_new'>;

export type Tab = {
    name: string;
    type: TabType;
    value: `${TabTypesForIndex}_${number}` | 'create_new';
}
