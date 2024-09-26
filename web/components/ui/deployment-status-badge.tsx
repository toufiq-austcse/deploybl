import { DEPLOYMENT_STATUS } from '@/lib/constant';
import { badgeVariants } from '@/components/ui/badge';
import * as React from 'react';
<<<<<<< Updated upstream
import Loading from '@/components/ui/loading';
=======
>>>>>>> Stashed changes

const DeploymentStatusBadge = ({ status }: { status: string }) => {
  if (status === DEPLOYMENT_STATUS.FAILED) {

    return (
      <div
<<<<<<< Updated upstream
        className={`${badgeVariants({ variant: 'destructive' })} capitalize gap-2`}
      >
        <div>
          {status}
        </div>

=======
        className={`${badgeVariants({ variant: 'destructive' })} capitalize`}
      >
        {status}
>>>>>>> Stashed changes
      </div>
    );
  } else if (status === DEPLOYMENT_STATUS.PULLING) {
    return (
      <div
<<<<<<< Updated upstream
        className={`${badgeVariants({ variant: 'default' })} capitalize gap-2`}
      >
        <div>
          {status}
        </div>
        <Loading className="bg-accent" />
=======
        className={`${badgeVariants({ variant: 'default' })} capitalize`}
      >
        {status}
>>>>>>> Stashed changes
      </div>
    );
  } else if (status === DEPLOYMENT_STATUS.BUILDING) {
    return (
      <div
<<<<<<< Updated upstream
        className={`${badgeVariants({ variant: 'default' })} capitalize gap-2`}
      >
        <div>
          {status}
        </div>
        <Loading className="bg-white" />
=======
        className={`${badgeVariants({ variant: 'outline' })} capitalize`}
      >
        {status}
>>>>>>> Stashed changes
      </div>
    );
  }

  return (
    <div
<<<<<<< Updated upstream
      className={`${badgeVariants({ variant: 'success' })} capitalize`}
=======
      className={`${badgeVariants({ variant: 'secondary' })} capitalize`}
>>>>>>> Stashed changes
    >
      {status}
    </div>
  );
};
export default DeploymentStatusBadge;