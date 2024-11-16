import { FaChevronDown, FaGithub, FaRegCopy } from 'react-icons/fa';
import * as React from 'react';
import moment from 'moment';
import DeploymentStatusBadge from '@/components/ui/deployment-status-badge';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
import { onCopyUrlClicked } from '@/lib/utils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';
import Link from 'next/link';

const DeploymentDetails = () => {
  const {
    deploymentDetails,
    updateLatestDeploymentStatus,
    restartDeploymentContext,
    rebuildAndDeployContext,
    stopDeploymentContext,
  } = useDeploymentContext();

  const onRestartClicked = async (deploymentId: string) => {
    let response = await restartDeploymentContext(deploymentId);
    if (response.error) {
      toast.error(response.error);
      return;
    }
    toast('Deployment restarting...');
  };
  const onRebuildAndDeployClicked = async (deploymentId: string) => {
    let response = await rebuildAndDeployContext(deploymentId);
    if (response.error) {
      toast.error(response.error);
      return;
    }
    toast('Deployment rebuilding and deploying...');
  };

  const onStopClicked = async (deploymentId: string) => {
    let response = await stopDeploymentContext(deploymentId);
    if (response.error) {
      toast.error(response.error);
      return;
    }
    toast('Deployment stopping...');
  };

  return (
    deploymentDetails && (
      <div>
        <div className="flex justify-between">
          <div className="flex gap-2">
            <p className="text-3xl">{deploymentDetails?.title}</p>
            <div className="flex flex-col justify-center">
              <DeploymentStatusBadge status={deploymentDetails.latest_status} />
            </div>
          </div>
          <div>
            <DropdownMenu>
              <DropdownMenuTrigger>
                <Button variant="secondary" className=" flex flex-row justify-around gap-2">
                  Actions
                  <FaChevronDown />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem onClick={() => onRestartClicked(deploymentDetails._id)}>Restart</DropdownMenuItem>
                <DropdownMenuItem onClick={() => onRebuildAndDeployClicked(deploymentDetails._id)}>
                  Rebuild & Redeploy
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => onStopClicked(deploymentDetails._id)}>Stop</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
        <Link href={deploymentDetails.repository_url} className="flex gap-2 max-w-fit" target="_blank">
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
            {deploymentDetails.domain_url && (
              <div className="flex flex-row gap-2 text-blue-500">
                <Link href={deploymentDetails.domain_url as string} target="_blank">
                  {deploymentDetails.domain_url}
                </Link>
                <div className="flex flex-col justify-center">
                  <FaRegCopy
                    className="cursor-pointer"
                    onClick={() => onCopyUrlClicked(deploymentDetails.domain_url)}
                  />
                </div>
              </div>
            )}
          </div>

          <div className="flex flex-row-reverse min-w-50">
            {deploymentDetails.last_deployed_at ? (
              <p>Last Deployed At : {moment(deploymentDetails.last_deployed_at).fromNow()}</p>
            ) : (
              <p>Not Deployed yet</p>
            )}
          </div>
        </div>
      </div>
    )
  );
};
export default DeploymentDetails;
