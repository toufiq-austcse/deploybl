import { useAuthContext } from '@/contexts/useAuthContext';
import React, { ComponentType, useEffect } from 'react';
import { useRouter } from 'next/navigation';


const PrivateRoute =<P extends object>(Component: ComponentType<P>) => {
  return (props: any) => {
    const { currentUser } = useAuthContext();
    const router = useRouter();
    // If user is not logged in, return login component

    useEffect(() => {
      if (!currentUser) {
        router.push("/login");
        return;
      }
    }, []);

    return <Component {...props} />;
  };
};

export default PrivateRoute;
