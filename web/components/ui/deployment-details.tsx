import Link from 'next/link';
import { FaGithub, FaRegCopy } from 'react-icons/fa';
import * as React from 'react';
import { useEffect } from 'react';
import moment from 'moment';
import DeploymentStatusBadge from '@/components/ui/deployment-status-badge';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
import { onCopyUrlClicked } from '@/lib/utils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';

const DeploymentDetails = () => {
  const { deploymentDetails, updateLatestDeploymentStatus } = useDeploymentContext();

  let interval: NodeJS.Timeout;

  useEffect(() => {
    interval = setInterval(() => {
      updateLatestDeploymentStatus(deploymentDetails._id);
    }, +(process.env.NEXT_PUBLIC_PULL_DELAY_MS as string));

    return () => clearInterval(interval);

  }, [deploymentDetails]);

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
        <Link className="flex flex-row gap-2" href={deploymentDetails.repository_url} target="_blank">
          <p className="underline">{deploymentDetails?.repository_name}</p>
          <p className="underline">{deploymentDetails?.branch_name}</p>
        </Link>
      </Link>

      <div className="flex gap-2 justify-between">
        <div className="min-w-50">
          {deploymentDetails.domain_url && <div className="flex flex-row gap-2 text-blue-500">
            <Link href={deploymentDetails.domain_url as string}
                  target="_blank">{deploymentDetails.domain_url}</Link>
            <div className="flex flex-col justify-center">
              <FaRegCopy className="cursor-pointer" onClick={() => onCopyUrlClicked(deploymentDetails.domain_url)} />
            </div>
          </div>}
        </div>
        <div>
          <DropdownMenu>
            <DropdownMenuTrigger className="rounded-md px-3 py-2 text-sm font-medium hover:bg-gray-700 hover:text-white underline hover:cursor-pointer">Actions</DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem>Restart</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Rebuild & Rdeploy</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

        </div>

        <div className="flex flex-row-reverse min-w-50">
          {deploymentDetails.last_deployed_at ?
            <p>Last Deployed At : {moment(deploymentDetails.last_deployed_at).fromNow()}</p> :
            <p>Not Deployed yet</p>}
        </div>
      </div>

    </div>
  );
};
export default DeploymentDetails;