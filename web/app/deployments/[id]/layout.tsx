'use client';
import '@/styles/globals.css';
import * as React from 'react';
import { Separator } from '@/components/ui/separator';
import DeploymentDetailsSidebar from '@/components/ui/deployment-details-sidebar';
import DeploymentDetails from '@/components/ui/deployment-details';
import { DeploymentContextProvider } from '@/contexts/useDeploymentContext';
import PrivateRoute from '@/components/private-route';

const DeploymentDetailsLayout = ({ children, params }: { children: React.ReactNode; params: { id: string } }) => {
  console.log('params', params);
  return (
    <div>
      <DeploymentContextProvider params={params}>
        <DeploymentDetails />
        <Separator orientation="horizontal" className="m-2" />
        <div className="flex h-screen">
          <DeploymentDetailsSidebar />
          <Separator orientation="vertical" />
          <div className="p-2 w-full">{children}</div>
        </div>
      </DeploymentContextProvider>
    </div>
  );
};
export default PrivateRoute(DeploymentDetailsLayout);
