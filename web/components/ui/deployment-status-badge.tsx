import { DEPLOYMENT_STATUS } from '@/lib/constant';
import { badgeVariants } from '@/components/ui/badge';
import * as React from 'react';
import Loading from '@/components/ui/loading';

const DeploymentStatusBadge = ({ status }: { status: string }) => {
  if (status === DEPLOYMENT_STATUS.FAILED) {

    return (
      <div
        className={`${badgeVariants({ variant: 'destructive' })} capitalize gap-2`}
      >
        <div>
          {status}
        </div>

      </div>
    );
  } else if (status === DEPLOYMENT_STATUS.PULLING) {
    return (
      <div
        className={`${badgeVariants({ variant: 'default' })} capitalize gap-2`}
      >
        <div>
          {status}
        </div>
        <Loading className="bg-accent" />
      </div>
    );
  } else if (status === DEPLOYMENT_STATUS.BUILDING) {
    return (
      <div
        className={`${badgeVariants({ variant: 'default' })} capitalize gap-2`}
      >
        <div>
          {status}
        </div>
        <Loading className="bg-white" />
      </div>
    );
  }

  return (
    <div
      className={`${badgeVariants({ variant: 'success' })} capitalize`}
    >
      {status}
    </div>
  );
};
export default DeploymentStatusBadge;