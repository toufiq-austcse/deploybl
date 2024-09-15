import Link from 'next/link';
import * as React from 'react';

const DeploymentDetailsNavbar = ({ deploymentId }: { deploymentId: string }) => {
  return (
    <nav className="flex flex-col w-1/6">
      <Link
        href={`/deployments/${deploymentId}/environments`}
        className="flex items-center gap-3 rounded-md px-3 py-2 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
      >
        <span className="hidden sm:block">Environments</span>
      </Link>

      <Link
        href={`/deployments/${deploymentId}/settings`}
        className="flex items-center gap-3 rounded-md px-3 py-2 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
        prefetch={false}
      >
        <span className="hidden sm:block">Settings</span>
      </Link>
    </nav>
  );
};
export default DeploymentDetailsNavbar;