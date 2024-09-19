'use client';
import { NextPage } from 'next';
import {
  Form,
  FormControl, FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form';
import { useForm } from 'react-hook-form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useState } from 'react';
import { MdDelete } from 'react-icons/md';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';

const EnvironmentPage: NextPage = () => {
  const { deploymentDetails } = useDeploymentContext();
  const [envs, setEnvs] = useState(deploymentDetails?.env);
  console.log('envs ', envs);
  const form = useForm();
  return (
    <div>
      <h1 className="font-bold text-2xl">Environment Variable</h1>
      <Form {...form}>
        <form id="login-form">
          <div className="flex flex-row gap-2 justify-center m-2">
            <h1 className="w-2/5">Key</h1>
            <h1 className="w-2/5">Value</h1>
          </div>
          {envs && Object.keys(envs as any).map((key, index, array) => {
            return <div key={index} className="flex flex-row gap-2 justify-center m-2">
              <div className="w-2/5">
                <FormField
                  control={form.control}
                  name="key"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <Input {...field} value={key} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-2/5">
                <FormField
                  control={form.control}
                  name="key"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <Input  {...field} value={envs[key]} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <div>
                <MdDelete onClick={() => {
                  envs.splice(index, 1);
                  console.log('newEnvs ', envs);
                  setEnvs(() => [...envs]);
                }} />
              </div>

            </div>;
          })}


        </form>
      </Form>
      <div className="flex flex-row-reverse gap-2">
        <Button
          onClick={() => {
            setEnvs(prev => [...prev, {}]);
          }}
          size="sm"
          className="my-2"
        >
          Save, Redeploy
        </Button>
        <Button
          onClick={() => {
            setEnvs(prev => [...prev, {}]);
          }}
          size="sm"
          className="my-2"
          variant="outline"
        >
          Add environment variable
        </Button>

      </div>

    </div>
  );
};
export default EnvironmentPage;