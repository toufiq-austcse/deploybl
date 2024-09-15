'use client';
import '@/styles/globals.css';
import * as React from 'react';
import { Separator } from '@/components/ui/separator';
import DeploymentDetailsNavbar from '@/components/ui/deployment-details-navbar';
import DeploymentDetails from '@/components/ui/deployment-details';


const DeploymentDetailsLayout = ({ children, params }: { children: React.ReactNode, params: { id: string } }) => {
  const deploymentId = params.id;

  return (
    <div>
      <DeploymentDetails />
      <Separator orientation="horizontal" className="m-2" />
      <div className="flex h-screen">
        <DeploymentDetailsNavbar deploymentId={deploymentId} />
        <Separator orientation="vertical" />
        <div className="p-2 w-full">
          {children}
        </div>
      </div>
    </div>
  );
};
export default DeploymentDetailsLayout;