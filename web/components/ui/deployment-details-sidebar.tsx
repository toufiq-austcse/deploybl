import Link from 'next/link';
import * as React from 'react';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';

const DeploymentDetailsSidebar = () => {
  const { deploymentDetails } = useDeploymentContext();
  return (
    <nav className="flex flex-col w-1/6">
      <Link
        href={`/deployments/${deploymentDetails?._id}/settings`}
        className="flex items-center gap-3 rounded-md px-3 py-2 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
        prefetch={false}
      >
        <h1 className="hidden sm:block">Settings</h1>
      </Link>
      <Link
        href={`/deployments/${deploymentDetails?._id}/environments`}
        className="flex items-center gap-3 rounded-md px-3 py-2 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
      >
        <h1 className="hidden sm:block">Environments</h1>
      </Link>
    </nav>
  );
};
export default DeploymentDetailsSidebar;
