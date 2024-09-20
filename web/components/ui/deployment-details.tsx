import Link from 'next/link';
import { FaGithub, FaRegCopy } from 'react-icons/fa';
import * as React from 'react';
import moment from 'moment';
import DeploymentStatusBadge from '@/components/ui/deployment-status-badge';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';

const DeploymentDetails = () => {
  const { deploymentDetails } = useDeploymentContext();
  console.log('deploymentDetails ', deploymentDetails);
  return (
    <>
      <div className="flex gap-2">
        <p className="text-3xl">{deploymentDetails?.title}</p>
        <DeploymentStatusBadge status={deploymentDetails?.latest_status as any} />
      </div>
      <Link href={''} className="flex gap-2">
        <div className="flex flex-col justify-center">
          <FaGithub />
        </div>
        <div className="flex flex-row gap-2">
          <p className="underline">{deploymentDetails?.repository_name}</p>
          <p className="underline">{deploymentDetails?.branch_name}</p>
        </div>
      </Link>

      <div className="flex gap-2 justify-between">
        <div className="min-w-50">
          {deploymentDetails?.domain_url && <div className="flex flex-row gap-2 text-blue-500">
            <Link href={''}>{deploymentDetails.domain_url}</Link>
            <div className="flex flex-col justify-center">
              <FaRegCopy />
            </div>
          </div>}
        </div>

        <div className="flex flex-row-reverse min-w-50">
          {deploymentDetails?.last_deployed_at ?
            <p>Last Deployed : {moment(deploymentDetails.last_deployed_at).fromNow()}</p> : <p>Not Deployed yet</p>}
        </div>
      </div>

    </>
  );
};
export default DeploymentDetails;