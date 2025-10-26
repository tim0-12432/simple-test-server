import { useEffect, useState } from 'react';
import { ThemeProvider } from "@/components/theme-provider"
import { Tabs, TabsContent, TabsList, TabsTrigger } from './components/ui/tabs'
import type { Tab } from './types/Tab';
import { TabFactory } from './tabs/TabFactory';
import { ModeToggle } from './components/mode-toggle';
import { getIconForTabType } from './lib/tabs';
import type { Container } from './types/Container';
import { request } from './lib/api';

function App() {
  const [tabs, setTabs] = useState<Tab[]>(addCreateTab([]));

  useEffect(() => {
    loadTabs();
  }, []);

  function loadTabs() {
    (async () => {
      try {
        const containers: Container[] = await request("GET", `/containers`);
        if (containers) {
          setTabs(addCreateTab(containers.map(container => ({
            name: container.name,
            id: container.container_id,
            type: container.type
          }))));
        } else {
          console.error("Failed to load containers.");
        }
      } catch (err) {
        console.error(`Error loading containers: ${err instanceof Error ? err.message : "Unknown error"}`);
      }
    })();
  }

  function addCreateTab(tabs: Tab[]) {
    const createTab: Tab = {
      name: "Create New",
      id: "create_new",
      type: "create_new"
    };
    if (!tabs.some(tab => tab.type === createTab.type)) {
      tabs.push(createTab);
    }
    return tabs;
  }

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Tabs defaultValue={tabs[0].id} className="w-full h-screen" key={tabs.length}>
        <nav>
          <TabsList className="w-full p-0 mt-2 bg-background justify-start border-b rounded-none">
            {tabs.map((tab) => (
              <TabsTrigger
                key={tab.id}
                value={tab.id}
                className="rounded-none bg-background h-full data-[state=active]:shadow-none border border-transparent border-b-border data-[state=active]:border-border data-[state=active]:border-b-background -mb-[2px] rounded-t cursor-pointer"
              >
                {getIconForTabType(tab.type)} <p className="text-[13px]">{tab.name}</p>
              </TabsTrigger>
            ))}
            <div className="flex-1" />
            <div className='h-full flex items-center justify-end pr-2'>
              <ModeToggle />
            </div>
          </TabsList>
        </nav>

        {tabs.map((tab) => (
          <TabsContent key={tab.id} value={tab.id} className='h-auto overflow-y-auto flex flex-col items-center'>
            <div className='w-full p-4'>
              {
                TabFactory(tab.type, {
                  id: tab.id,
                  type: tab.type,
                  reloadTabs: loadTabs
                })
              }
            </div>
            <div className='flex-1' />
            <footer className='w-full max-w-xs xl:max-w-md 2xl:max-w-lg'>
              <div className='p-2 border-t border-border text-center text-xs text-muted-foreground'>
                Created with <span className="animate-pulse">❤️</span> by <a className='underline-offset-4 hover:underline' href="https://github.com/tim0-12432">Tim0-12432</a>
              </div>
            </footer>
          </TabsContent>
        ))}
      </Tabs>
    </ThemeProvider>
  )
}

export default App
