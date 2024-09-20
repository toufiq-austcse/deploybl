'use client';
import '@/styles/globals.css';
import * as React from 'react';
import { Separator } from '@/components/ui/separator';
import DeploymentDetailsNavbar from '@/components/ui/deployment-details-navbar';
import DeploymentDetails from '@/components/ui/deployment-details';
import { useHttpClient } from '@/api/http/useHttpClient';
import { useEffect, useState } from 'react';
import { DeploymentDetailsType } from '@/api/http/types/deployment_type';


const DeploymentDetailsLayout = ({ children, params }: { children: React.ReactNode, params: { id: string } }) => {
  const [deploymentDetails, setDeploymentDetails] = useState<DeploymentDetailsType>();
  let { getDeploymentDetails, loading } = useHttpClient();

  const deploymentId = params.id;

  useEffect(() => {
    getDeploymentDetails(deploymentId).then(data => {
      setDeploymentDetails(() => data);
      console.log('Loading ', loading);
    }).catch(err => {
      console.log(err);
    });
  }, [deploymentId]);

  return (
    <div>
      {loading && <div>Loading</div>}
      {!loading && deploymentDetails && <div>
        <DeploymentDetails deploymentDetails={deploymentDetails as any} />
        <Separator orientation="horizontal" className="m-2" />
        <div className="flex h-screen">
          <DeploymentDetailsNavbar deploymentId={deploymentId} />
          <Separator orientation="vertical" />
          <div className="p-2 w-full">
            {children}
          </div>
        </div>
      </div>}

    </div>
  );
};
export default DeploymentDetailsLayout;