import Link from 'next/link';
import { FaGithub, FaRegCopy } from 'react-icons/fa';
import * as React from 'react';
import { useEffect } from 'react';
import moment from 'moment';
import DeploymentStatusBadge from '@/components/ui/deployment-status-badge';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
import { toast } from 'sonner';
import { DEPLOYMENT_STATUS } from '@/lib/constant';

const DeploymentDetails = () => {
  const { deploymentDetails, updateLatestDeploymentStatus } = useDeploymentContext();

  let interval: NodeJS.Timeout;

  useEffect(() => {
    console.log('deployment details changed');

    if (deploymentDetails && deploymentDetails.latest_status != DEPLOYMENT_STATUS.LIVE && deploymentDetails.latest_status != DEPLOYMENT_STATUS.FAILED) {
      interval = setInterval(() => {
        updateLatestDeploymentStatus(deploymentDetails._id);
      }, +(process.env.NEXT_PUBLIC_PULL_DELAY_MS as string));

      return () => clearInterval(interval);
    }

  }, [deploymentDetails]);

  const onCopyUrlClicked = async () => {
    await navigator.clipboard.writeText(deploymentDetails?.domain_url as string);
    toast('Copied to clipboard');
  };

  return (
    deploymentDetails &&
    <div>
      <div className="flex gap-2">
        <p className="text-3xl">{deploymentDetails?.title}</p>
        <DeploymentStatusBadge status={deploymentDetails.latest_status} />
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
          {deploymentDetails.domain_url && <div className="flex flex-row gap-2 text-blue-500">
            <Link href={deploymentDetails.domain_url as string}
                  target="_blank">{deploymentDetails.domain_url}</Link>
            <div className="flex flex-col justify-center">
              <FaRegCopy className="cursor-pointer" onClick={onCopyUrlClicked} />
            </div>
          </div>}
        </div>

        <div className="flex flex-row-reverse min-w-50">
          {deploymentDetails.last_deployed_at ?
            <p>Last Deployed : {moment(deploymentDetails.last_deployed_at).fromNow()}</p> :
            <p>Not Deployed yet</p>}
        </div>
      </div>

    </div>
  );
};
export default DeploymentDetails;