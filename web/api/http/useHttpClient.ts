import React from 'react';
import axios, { AxiosError } from 'axios';
import FormData from 'form-data';
import {
  DeploymentDetailsType,
  DeploymentLatestStatusType,
  DeploymentType,
  RepoDetailsType
} from '@/api/http/types/deployment_type';

export function useHttpClient() {
  const [loading, setLoading] = React.useState(false);

  const getDeploymentDetails = async (deploymentId: string): Promise<DeploymentDetailsType> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/deployments/${deploymentId}`;

      const response = await axios.get(url);
      return response.data.data;
    } catch (err) {
      let message = (err as any).message;
      throw new Error(message);
    } finally {
      setLoading(false);
    }
  };

  const listDeployments = async (page: number, limit: number): Promise<{
    data: DeploymentType[] | null;
    error: string | null;
  }> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/deployments`;

      const response = await axios.get(url);
      return {
        data: response.data.data,
        error: null
      };
    } catch (err) {
      return handleError(err);
    } finally {
      setLoading(false);
    }
  };

  const handleError = (err: any) => {
    if (axios.isAxiosError(err)) {
      let error = '';
      let errorResponse: any = (err as AxiosError).response?.data;
      if (errorResponse) {
        error = errorResponse.errors.join(',');
      } else {
        error = (err as AxiosError).message;
      }
      return { data: null, error };
    }
    let message = (err as any).message;
    return { data: null, error: message };
  };

  const getRepoDetails = async (repoUrl: string): Promise<{
    data: RepoDetailsType | null;
    error: string | null;
  }> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/repositories?repo_url=${repoUrl}`;

      const response = await axios.get(url);
      return {
        data: response.data.data,
        error: null
      };
    } catch (err) {
      return handleError(err);
    } finally {
      setLoading(false);
    }
  };

  const createDeployment = async (body: {
    title: string
    branch_name: string,
    docker_file_path: string,
    env: object,
    repository_url: string,
    root_dir: string | null
  }): Promise<{
    data: DeploymentType | null;
    error: string | null;
  }> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/deployments`;
      const response = await axios.post(url, body);
      return {
        data: response.data.data,
        error: null
      };
    } catch (err) {
      return handleError(err);
    } finally {
      setLoading(false);
    }
  };

  const getDeploymentLatestStatus = async (deploymentIds: string[]): Promise<{
    data: DeploymentLatestStatusType[] | null;
    error: string | null;
  }> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/deployments/latest-status?ids=${deploymentIds}`;
      const response = await axios.get(url);
      return {
        data: response.data.data,
        error: null
      };
    } catch (err) {
      return handleError(err);
    } finally {
      setLoading(false);
    }
  };


  return {
    getDeploymentDetails,
    listDeployments,
    getRepoDetails,
    createDeployment,
    getDeploymentLatestStatus,
    loading
  };
};