import React from 'react';
import axios from 'axios';
import FormData from 'form-data';
import { DeploymentDetailsType, DeploymentType } from '@/api/http/types/deployment_type';

export function useHttpClient() {
  const [loading, setLoading] = React.useState(false);

  const GetDeploymentDetails = async (deploymentId: string): Promise<DeploymentDetailsType> => {
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

  const listDeployments = async (page: number, limit: number): Promise<DeploymentType[]> => {
    setLoading(true);
    try {
      let url = `${process.env.NEXT_PUBLIC_JUST_DEPLOY_API_URL}/deployments`;

      const response = await axios.get(url);
      return response.data.data;
    } catch (err) {
      let message = (err as any).message;
      throw new Error(message);
    } finally {
      setLoading(false);
    }
  };


  return {
    GetDeploymentDetails,
    listDeployments,
    loading
  };
}