import React from 'react';
import { ThemeProvider } from "@/components/theme-provider"
import { Tabs, TabsContent, TabsList, TabsTrigger } from './components/ui/tabs'
import type { Tab } from './types/Tab';
import TabFactory from './tabs/TabFactory';
import { ModeToggle } from './components/mode-toggle';
import { getIconForTabType } from './lib/tabs';

function App() {

    const tabs: Tab[] = [
    {
      name: "Best MQTT",
      type: "MQTT",
      value: "MQTT_0"
    },
    {
      name: "FTP Safe",
      type: "FTP",
      value: "FTP_0"
    },
    {
      name: "Web App",
      type: "Web",
      value: "Web_0"
    },
    {
      name: "",
      type: "create_new",
      value: "create_new"
    }
  ];

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Tabs defaultValue={tabs[0].value} className="w-full h-screen">
        <nav>
          <TabsList className="w-full p-0 mt-2 bg-background justify-start border-b rounded-none">
            {tabs.map((tab) => (
              <TabsTrigger
                key={tab.value}
                value={tab.value}
                className="rounded-none bg-background h-full data-[state=active]:shadow-none border border-transparent border-b-border data-[state=active]:border-border data-[state=active]:border-b-background -mb-[2px] rounded-t"
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
          <TabsContent key={tab.value} value={tab.value} className='h-auto overflow-y-auto flex flex-col items-center'>
            <div className='w-full p-4'>
              {
                TabFactory(tab.type, {
                  id: tab.value,
                  type: tab.type
                })
              }
            </div>
            <div className='flex-1' />
            <footer className='w-full max-w-xs xl:max-w-md 2xl:max-w-lg'>
              <div className='p-2 border-t border-border text-center text-xs text-muted-foreground'>
                Created with ❤️ by <a className='underline-offset-4 hover:underline' href="https://github.com/tim0-12432">Tim0-12432</a>
              </div>
            </footer>
          </TabsContent>
        ))}
      </Tabs>
    </ThemeProvider>
  )
}

export default App
