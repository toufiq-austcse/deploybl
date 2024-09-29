import { DeploymentDetailsType, UpdateDeploymentReqBody } from '@/api/http/types/deployment_type';
import React, { useContext, useEffect } from 'react';
import { useHttpClient } from '@/api/http/useHttpClient';

type DeploymentContextType = {
  deploymentDetails: DeploymentDetailsType | null;
  updateDeploymentDetails: (deploymentId: string, body: UpdateDeploymentReqBody) => Promise<{ error: string | null }>;
  updateLatestDeploymentStatus: (deploymentId: string) => Promise<{ error: string | null }>;
  updateEnv: (deploymentId: string, env: object) => Promise<{ error: string | null }>;
  loading: boolean;
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
  const { getDeploymentDetails, updateDeployment, getDeploymentLatestStatus, updateDeploymentEnv } = useHttpClient();
  const [loading, setLoading] = React.useState(true);

  useEffect(() => {
    getDeploymentDetails(params.id).then(data => {
      setDeploymentDetails(() => data);
    }).catch(err => {
      console.log(err);
    }).finally(() => {
      setLoading(false);
    });


  }, []);

  const updateDeploymentDetails = async (deploymentId: string, body: UpdateDeploymentReqBody): Promise<{
    error: string | null
  }> => {
    let response = await updateDeployment(deploymentId, body);
    if (response.error) {
      return {
        error: response.error
      };
    }
    console.log(response);
    setDeploymentDetails((prevState: any) => {
      return {
        ...prevState,
        ...response.data
      };
    });
    return {
      error: null
    };
  };

  const updateLatestDeploymentStatus = async (deploymentId: string): Promise<{
    error: string | null
  }> => {
    let response = await getDeploymentLatestStatus([deploymentId]);
    if (response.error) {
      return {
        error: response.error
      };
    }
    setDeploymentDetails((prevState: any) => {
      return {
        ...prevState,
        ...response.data?.[0]
      };
    });
    return {
      error: null

    };
  };

  const updateEnv = async (deploymentId: string, env: object): Promise<{
    error: string | null
  }> => {
    let response = await updateDeploymentEnv(deploymentId, env);
    if (response.error) {
      return {
        error: response.error
      };
    }
    setDeploymentDetails((prevState: any) => {
      return {
        ...prevState,
        ...response.data
      };
    });
    return {
      error: null
    };
  };

  const value: DeploymentContextType = {
    deploymentDetails,
    updateDeploymentDetails,
    updateLatestDeploymentStatus,
    updateEnv,
    loading
  };
  return (
    <DeploymentContext.Provider value={value}>
      {loading ? <div>Loading...</div> : children}
    </DeploymentContext.Provider>
  );
};