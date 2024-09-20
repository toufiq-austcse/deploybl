import { DEPLOYMENT_STATUS } from '@/lib/constant';
import { badgeVariants } from '@/components/ui/badge';
import * as React from 'react';

const DeploymentStatusBadge = ({ status }: { status: string }) => {
  if (status === DEPLOYMENT_STATUS.FAILED) {

    return (
      <div
        className={`${badgeVariants({ variant: 'destructive' })} capitalize`}
      >
        {status}
      </div>
    );
  } else if (status === DEPLOYMENT_STATUS.PULLING) {
    return (
      <div
        className={`${badgeVariants({ variant: 'default' })} capitalize`}
      >
        {status}
      </div>
    );
  } else if (status === DEPLOYMENT_STATUS.BUILDING) {
    return (
      <div
        className={`${badgeVariants({ variant: 'outline' })} capitalize`}
      >
        {status}
      </div>
    );
  }

  return (
    <div
      className={`${badgeVariants({ variant: 'secondary' })} capitalize`}
    >
      {status}
    </div>
  );
};
export default DeploymentStatusBadge;