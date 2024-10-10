import { DeploymentDetailsType, UpdateDeploymentReqBody } from '@/api/http/types/deployment_type';
import React, { useContext, useEffect } from 'react';
import { useHttpClient } from '@/api/http/useHttpClient';
import { toast } from 'sonner';

type DeploymentContextType = {
  deploymentDetails: DeploymentDetailsType | null;
  updateDeploymentDetails: (deploymentId: string, body: UpdateDeploymentReqBody) => Promise<{ error: string | null }>;
  updateLatestDeploymentStatus: (deploymentId: string) => Promise<{ error: string | null }>;
  updateEnv: (deploymentId: string, env: object) => Promise<{ error: string | null }>;
  restartDeploymentContext: (deploymentId: string) => Promise<{ error: string | null }>;
  rebuildAndDeployContext: (deploymentId: string) => Promise<{ error: string | null }>;
  stopDeploymentContext: (deploymentId: string) => Promise<{ error: string | null }>;
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
  const {
    getDeploymentDetails,
    updateDeployment,
    getDeploymentLatestStatus,
    updateDeploymentEnv,
    restartDeployment,
    rebuildAndDeploy,
    stopDeployment
  } = useHttpClient();
  const [loading, setLoading] = React.useState(true);

  useEffect(() => {
    getDeploymentDetails(params.id).then(response => {
      if (!response.isSuccessful && response.code !== 401) {
        toast.error(response.error);
        return;
      }
      setDeploymentDetails(response.data);
    }).finally(() => {
      setLoading(false);
    });
  }, []);

  const updateDeploymentDetails = async (deploymentId: string, body: UpdateDeploymentReqBody): Promise<{
    error: string | null
  }> => {
    let response = await updateDeployment(deploymentId, body);
    if (!response.isSuccessful && response.code !== 401) {
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
    if (!response.isSuccessful && response.code !== 401) {
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
    if (!response.isSuccessful && response.code !== 401) {
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
  const restartDeploymentContext = async (deploymentId: string) => {
    let response = await restartDeployment(deploymentId);
    if (!response.isSuccessful && response.code !== 401) {
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
  const rebuildAndDeployContext = async (deploymentId: string) => {
    let response = await rebuildAndDeploy(deploymentId);
    if (!response.isSuccessful && response.code !== 401) {
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
  const stopDeploymentContext = async (deploymentId: string) => {
    let response = await stopDeployment(deploymentId);
    if (!response.isSuccessful && response.code !== 401) {
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
    loading,
    restartDeploymentContext,
    rebuildAndDeployContext,
    stopDeploymentContext
  };
  return (
    <DeploymentContext.Provider value={value}>
      {loading ? <div>Loading...</div> : children}
    </DeploymentContext.Provider>
  );
};