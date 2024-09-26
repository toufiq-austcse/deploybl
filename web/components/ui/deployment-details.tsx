import Link from 'next/link';
import { FaGithub, FaRegCopy } from 'react-icons/fa';
import * as React from 'react';
import moment from 'moment';
import DeploymentStatusBadge from '@/components/ui/deployment-status-badge';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
<<<<<<< Updated upstream
import { useEffect } from 'react';
import { useHttpClient } from '@/api/http/useHttpClient';
import { DEPLOYMENT_STATUS } from '@/lib/constant';
import { toast } from 'sonner';

const DeploymentDetails = () => {
  const { loading, getDeploymentLatestStatus } = useHttpClient();
  const { deploymentDetails } = useDeploymentContext();
  const [latestDeploymentDetails, setLatestDeploymentDetails] = React.useState(deploymentDetails);
  let interval: NodeJS.Timeout;
  useEffect(() => {
    if (latestDeploymentDetails && latestDeploymentDetails.latest_status != DEPLOYMENT_STATUS.LIVE && latestDeploymentDetails.latest_status != DEPLOYMENT_STATUS.FAILED) {
      interval = setInterval(() => {
        getDeploymentLatestStatus([latestDeploymentDetails?._id]).then(response => {
          if (response.error) {
            console.log(response.error);
          } else {
            console.log(response);
            // @ts-ignore
            if (response.data?.length > 0) {
              console.log('setting...');
              setLatestDeploymentDetails((prevState: any) => {
                // @ts-ignore
                return {
                  ...prevState,
                  latest_status: (response as any).data[0].latest_status,
                  last_deployed_at: (response as any).data[0].last_deployed_at,
                  domain_url: (response as any).data[0].domain_url
                };
              });
            }
          }
        });
      }, +(process.env.NEXT_PUBLIC_PULL_DELAY_MS as string));

      return () => clearInterval(interval);
    }

  }, [latestDeploymentDetails]);

  const onCopyUrlClicked = async () => {
    await navigator.clipboard.writeText(latestDeploymentDetails?.domain_url as string);
    toast('Copied to clipboard');
  };

  return (
    <>
      <div className="flex gap-2">
        <p className="text-3xl">{latestDeploymentDetails?.title}</p>
        <DeploymentStatusBadge status={latestDeploymentDetails?.latest_status as any} />
=======

const DeploymentDetails = () => {
  const { deploymentDetails } = useDeploymentContext();
  console.log('deploymentDetails ', deploymentDetails);
  return (
    <>
      <div className="flex gap-2">
        <p className="text-3xl">{deploymentDetails?.title}</p>
        <DeploymentStatusBadge status={deploymentDetails?.latest_status as any} />
>>>>>>> Stashed changes
      </div>
      <Link href={''} className="flex gap-2">
        <div className="flex flex-col justify-center">
          <FaGithub />
        </div>
        <div className="flex flex-row gap-2">
<<<<<<< Updated upstream
          <p className="underline">{latestDeploymentDetails?.repository_name}</p>
          <p className="underline">{latestDeploymentDetails?.branch_name}</p>
=======
          <p className="underline">{deploymentDetails?.repository_name}</p>
          <p className="underline">{deploymentDetails?.branch_name}</p>
>>>>>>> Stashed changes
        </div>
      </Link>

      <div className="flex gap-2 justify-between">
        <div className="min-w-50">
<<<<<<< Updated upstream
          {latestDeploymentDetails?.domain_url && <div className="flex flex-row gap-2 text-blue-500">
            <Link href={latestDeploymentDetails?.domain_url as string}
                  target="_blank">{latestDeploymentDetails?.domain_url}</Link>
            <div className="flex flex-col justify-center">
              <FaRegCopy className="cursor-pointer" onClick={onCopyUrlClicked} />
=======
          {deploymentDetails?.domain_url && <div className="flex flex-row gap-2 text-blue-500">
            <Link href={''}>{deploymentDetails.domain_url}</Link>
            <div className="flex flex-col justify-center">
              <FaRegCopy />
>>>>>>> Stashed changes
            </div>
          </div>}
        </div>

        <div className="flex flex-row-reverse min-w-50">
<<<<<<< Updated upstream
          {latestDeploymentDetails?.last_deployed_at ?
            <p>Last Deployed : {moment(latestDeploymentDetails?.last_deployed_at).fromNow()}</p> :
            <p>Not Deployed yet</p>}
=======
          {deploymentDetails?.last_deployed_at ?
            <p>Last Deployed : {moment(deploymentDetails.last_deployed_at).fromNow()}</p> : <p>Not Deployed yet</p>}
>>>>>>> Stashed changes
        </div>
      </div>

    </>
  );
};
export default DeploymentDetails;