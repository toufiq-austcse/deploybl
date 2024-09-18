import { badgeVariants } from '@/components/ui/badge';
import Link from 'next/link';
import { FaGithub, FaRegCopy } from 'react-icons/fa';
import * as React from 'react';
import { DeploymentDetailsType } from '@/api/http/types/deployment_type';

const DeploymentDetails = ({ deploymentDetails }: { deploymentDetails: DeploymentDetailsType }) => {
  return (
    <>
      <div className="flex gap-2">
        <p className="text-3xl">{deploymentDetails.title}</p>
        <div className={`${badgeVariants({ variant: 'secondary' })} capitalize`}>Live</div>
      </div>
      <Link href={''} className="flex gap-2">
        <div className="flex flex-col justify-center">
          <FaGithub />
        </div>
        <div className="flex flex-row gap-2">
          <p className="underline">toufiq-austcse/test</p>
          <p className="underline">{deploymentDetails.branch_name}</p>
        </div>
      </Link>

      <div className="flex gap-2 justify-between ">
        <div className="flex flex-row gap-2 text-blue-500">
          <Link href={''}>https://test.com</Link>
          <div className="flex flex-col justify-center">
            <FaRegCopy />
          </div>
        </div>
        <div className="flex flex-row-reverse">
          <p>Last Deployed At : 2 hours ago</p>
        </div>
      </div>

    </>
  );
};
export default DeploymentDetails;