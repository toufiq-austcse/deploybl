import React from 'react';
import axios, { AxiosError } from 'axios';
import {
  DeploymentDetailsType,
  DeploymentEventType,
  DeploymentLatestStatusType,
  DeploymentType,
  RepoDetailsType,
  TPaginationResponse,
  TResponse,
  UpdateDeploymentReqBody
} from '@/api/http/types/deployment_type';
import { useAuthContext } from '@/contexts/useAuthContext';

export function useHttpClient() {
  const [loading, setLoading] = React.useState(false);
  const { currentUser, logout } = useAuthContext();

  const getDeploymentDetails = async (deploymentId: string): Promise<TResponse<DeploymentDetailsType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/${deploymentId}`;

      const response = await axios.get(url, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null

      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }
  };

  const listDeployments = async (page: number, limit: number): Promise<TPaginationResponse<DeploymentType[]>> => {
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments?page=${page}&limit=${limit}`;
      const response = await axios.get(url, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        data: response.data.data,
        pagination: response.data.pagination,
        error: null,
        code: response.status
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: [],
        error,
        code,
        pagination: null
      };
    }
  };

  const handleError = async (err: any): Promise<{ error: string, code: number }> => {
    let code = 500;
    let error: string;
    if (axios.isAxiosError(err)) {
      let errorResponse: any = (err as AxiosError).response?.data;
      if (errorResponse) {
        if (typeof errorResponse === 'string') {
          error = errorResponse;
        } else {
          error = errorResponse.errors.join(',');
        }

        code = errorResponse.code;
      } else {
        error = (err as AxiosError).message;
      }

    } else {
      error = (err as any).message;
    }
    if (code === 401) {
      await logout();
    }

    return { error, code };
  };

  const getRepoDetails = async (repoUrl: string): Promise<TResponse<RepoDetailsType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/repositories?repo_url=${repoUrl}`;

      const response = await axios.get(url, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }
  };
  const getRepoBranches = async (repoUrl: string): Promise<TResponse<string[]>> => {
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/repositories/branches?repo_url=${repoUrl}`;

      const response = await axios.get(url, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data.map((branch: any) => branch.name),
        error: null
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: [],
        error,
        code
      };
    }
  };

  const createDeployment = async (body: {
    title: string
    branch_name: string,
    docker_file_path: string,
    env: object,
    repository_url: string,
    root_dir: string | null
  }): Promise<TResponse<DeploymentType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments`;
      const response = await axios.post(url, body, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }
  };

  const getDeploymentLatestStatus = async (deploymentIds: string[]): Promise<TResponse<DeploymentLatestStatusType[]>> => {
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/latest-status?ids=${deploymentIds}`;
      const response = await axios.get(url, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    }
  };

  const updateDeployment = async (deploymentId: string, body: UpdateDeploymentReqBody): Promise<TResponse<DeploymentDetailsType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/${deploymentId}`;
      const response = await axios.patch(url, body, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { error, code } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }

  };
  const updateDeploymentEnv = async (deploymentId: string, env: object): Promise<TResponse<DeploymentDetailsType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/${deploymentId}/env`;
      const response = await axios.patch(url, env, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { error, code } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }
  };
  const restartDeployment = async (deploymentId: string): Promise<TResponse<DeploymentDetailsType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/${deploymentId}/restart`;
      const response = await axios.post(url, {}, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { error, code } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }
  };
  const rebuildAndDeploy = async (deploymentId: string): Promise<TResponse<DeploymentDetailsType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/${deploymentId}/rebuild-and-redeploy`;
      const response = await axios.post(url, {}, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }
  };
  const stopDeployment = async (deploymentId: string): Promise<TResponse<DeploymentDetailsType>> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/${deploymentId}/stop`;
      const response = await axios.post(url, {}, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        error: null
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: null,
        error,
        code
      };
    } finally {
      setLoading(false);
    }
  };
  const listDeploymentEvents = async (deploymentId: string, page: number, limit: number): Promise<TPaginationResponse<DeploymentEventType[]>> => {
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/api/v1/deployments/${deploymentId}/events?page=${page}&limit=${limit}`;
      const response = await axios.get(url, {
        headers: {
          Authorization: `Bearer ${currentUser?.accessToken}`
        }
      });
      return {
        isSuccessful: true,
        code: response.status,
        data: response.data.data,
        pagination: response.data.pagination,
        error: null
      };
    } catch (err) {
      let { code, error } = await handleError(err);
      return {
        isSuccessful: false,
        data: [],
        error,
        code,
        pagination: null
      };
    }


  };


  return {
    getDeploymentDetails,
    listDeployments,
    getRepoDetails,
    createDeployment,
    getDeploymentLatestStatus,
    updateDeployment,
    updateDeploymentEnv,
    rebuildAndDeploy,
    restartDeployment,
    stopDeployment,
    getRepoBranches,
    listDeploymentEvents,
    loading
  };
}