import { useAuthContext } from '@/contexts/useAuthContext';
import { ComponentType, useEffect } from 'react';
import { useRouter } from 'next/navigation';

// Function that takes a component and returns a new component
const PublicRoute = <P extends object>(Component: ComponentType<P>) => {
  // Return a new functional component
  return (props: P) => {
    const { currentUser } = useAuthContext();
    const router = useRouter();

    useEffect(() => {
      if (currentUser) {
        router.push('/');
      }
    }, [currentUser, router]);

    // If user is not authenticated, render the component
    return <Component {...props} />;
  };
};

export default PublicRoute;