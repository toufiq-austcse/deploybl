import { DEPLOYMENT_STATUS } from '@/lib/constant';
import { badgeVariants } from '@/components/ui/badge';
import * as React from 'react';
import Loading from '@/components/ui/loading';

const DeploymentStatusBadge = ({ status }: { status: string }) => {
  if (status === DEPLOYMENT_STATUS.FAILED) {
    return <div className={`${badgeVariants({ variant: 'destructive' })} capitalize`}>{status}</div>;
  }
  if (
    status === DEPLOYMENT_STATUS.PULLING ||
    status === DEPLOYMENT_STATUS.BUILDING ||
    status === DEPLOYMENT_STATUS.QUEUED ||
    status === DEPLOYMENT_STATUS.DEPLOYING
  ) {
    return (
      <div className={`${badgeVariants({ variant: 'default' })} capitalize gap-2`}>
        <div>{status}</div>
        <Loading className="bg-accent" />
      </div>
    );
  }
  if (status === DEPLOYMENT_STATUS.LIVE) {
    return <div className={`${badgeVariants({ variant: 'success' })} capitalize`}>{status}</div>;
  }
  return <div className={`${badgeVariants({ variant: 'secondary' })} capitalize`}>{status}</div>;
};
export default DeploymentStatusBadge;
