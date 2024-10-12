'use client';
import { NextPage } from 'next';
import { Button } from '@/components/ui/button';
import * as React from 'react';
import { useEffect, useState } from 'react';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
import EnvironmentComponent, { EnvironmentVariableType } from '@/components/ui/environment-component';
import { convertEnvToObj, convertObjToEnv } from '@/lib/utils';
import ErrorAlert from '@/components/ui/error-alert';

const EnvironmentPage: NextPage = () => {
  const { deploymentDetails, updateEnv } = useDeploymentContext();
  const [envs, setEnvs] = useState<EnvironmentVariableType[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    if (deploymentDetails) {
      let envsArray = convertObjToEnv(deploymentDetails.env);
      setEnvs((prevState) => [...prevState, ...envsArray]);
    }
  }, []);
  const handleSaveAndRedeploy = async () => {
    setLoading(true);
    setError(null);
    let envObj = convertEnvToObj(envs);
    updateEnv(deploymentDetails?._id, envObj)
      .then((response) => {
        if (response.error) {
          setError(response.error);
        }
      })
      .finally(() => {
        setLoading(false);
      });
  };
  return (
    <div>
      <h1 className="font-bold text-2xl">Environment Variable</h1>
      {error && <ErrorAlert error={error} />}
      <div className="border 2px my-2">
        <EnvironmentComponent envs={envs} setEnvs={setEnvs} />
      </div>
      <div className="flex flex-row-reverse gap-2">
        <Button disabled={loading || envs.length === 0} onClick={handleSaveAndRedeploy} size="sm" className="my-2">
          Save & Redeploy
        </Button>
      </div>
    </div>
  );
};
export default EnvironmentPage;
