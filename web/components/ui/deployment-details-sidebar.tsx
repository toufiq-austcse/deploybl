'use client';
import Link from 'next/link';
import * as React from 'react';
import { useEffect } from 'react';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
import { usePathname } from 'next/navigation';

const DeploymentDetailsSidebar = () => {
  const { deploymentDetails } = useDeploymentContext();
  const pathname = usePathname();
  const [activeTab, setActiveTab] = React.useState('events');

  useEffect(() => {
    if (pathname.includes('settings')) {
      setActiveTab('settings');
    } else if (pathname.includes('environments')) {
      setActiveTab('environments');
    } else {
      setActiveTab('events');
    }
  }, [pathname]);
  return (
    <nav className="flex flex-col w-1/6">
      <Link
        href={`/deployments/${deploymentDetails?._id}/events`}
        className={`flex items-center gap-3 rounded-md px-3 py-2 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground ${
          activeTab === 'events' ? 'bg-accent text-accent-foreground' : ''
        }`}
        prefetch={false}
      >
        <h1 className="hidden sm:block">Events</h1>
      </Link>
      <Link
        href={`/deployments/${deploymentDetails?._id}/settings`}
        className={`flex items-center gap-3 rounded-md px-3 py-2 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground ${
          activeTab === 'settings' ? 'bg-accent text-accent-foreground' : ''
        }`}
        prefetch={false}
      >
        <h1 className="hidden sm:block">Settings</h1>
      </Link>
      <Link
        href={`/deployments/${deploymentDetails?._id}/environments`}
        className={`flex items-center gap-3 rounded-md px-3 py-2 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground ${
          activeTab === 'environments' ? 'bg-accent text-accent-foreground' : ''
        }`}
      >
        <h1 className="hidden sm:block">Environments</h1>
      </Link>
    </nav>
  );
};
export default DeploymentDetailsSidebar;
