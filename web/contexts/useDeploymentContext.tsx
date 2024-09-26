import { DeploymentDetailsType } from '@/api/http/types/deployment_type';
import React, { useContext, useEffect } from 'react';
import { useHttpClient } from '@/api/http/useHttpClient';

type DeploymentContextType = {
  deploymentDetails: DeploymentDetailsType | null;
}
const DeploymentContext = React.createContext({} as DeploymentContextType);

type DeploymentContextProviderProps = {
  children: React.ReactNode;
  params: { id: string }
}
export const useDeploymentContext = () => {
  return useContext(DeploymentContext);
};

export const DeploymentContextProvider = ({ children, params }: DeploymentContextProviderProps) => {
  const [deploymentDetails, setDeploymentDetails] = React.useState<DeploymentDetailsType | null>(null);
<<<<<<< Updated upstream
  const { loading, getDeploymentDetails, getDeploymentLatestStatus } = useHttpClient();

  useEffect(() => {
    getDeploymentDetails(params.id).then(data => {
=======
  const { loading, GetDeploymentDetails } = useHttpClient();

  useEffect(() => {
    GetDeploymentDetails(params.id).then(data => {
>>>>>>> Stashed changes
      setDeploymentDetails(() => data);
    }).catch(err => {
      console.log(err);
    });
<<<<<<< Updated upstream


=======
>>>>>>> Stashed changes
  }, []);

  const value: DeploymentContextType = {
    deploymentDetails: deploymentDetails
  };
  return (
    <DeploymentContext.Provider value={value}>
      {!loading && children}
    </DeploymentContext.Provider>
  );
};