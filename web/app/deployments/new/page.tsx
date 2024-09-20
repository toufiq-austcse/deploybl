'use client';
import { NextPage } from 'next';
import { Input } from '@/components/ui/input';
import EnvironmentComponent, { EnvironmentVariableType } from '@/components/ui/environment-component';
import { useState } from 'react';
import { RepoDetailsType } from '@/api/http/types/deployment_type';
import { z } from 'zod';
import * as React from 'react';
import { useHttpClient } from '@/api/http/useHttpClient';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import ErrorAlert from '@/components/ui/error-alert';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Button } from '@/components/ui/button';
import { useRouter } from 'next/navigation';

const validateRepositorySchema = z.object({
  repository_url: z
    .string()
    .url({
      message: 'Please enter valid url'
    })
});
const createDeploymentSchema = z.object({
  name: z.string({
    required_error: 'Name is required'
  }),
  repository_url: z.string({
    required_error: 'Repository URL is required'
  }).url({
    message: 'Please enter valid url'
  }),
  branch_name: z.string({
    required_error: 'Branch name is required'
  }),
  root_directory: z.string({
    required_error: 'Root directory is required'
  }),
  docker_file_path: z.string({
    required_error: 'Docker file path is required'

  })
});

const NewDeploymentPage: NextPage = () => {
  const router = useRouter();
  const { loading, createDeployment, getRepoDetails } = useHttpClient();
  const [repoDetails, setRepoDetails] = useState<RepoDetailsType | null>(null);
  const [error, setError] = React.useState<string | null>(null);
  const [envs, setEnvs] = useState<EnvironmentVariableType[]>([{ key: '', value: '' }]);

  const validateRepositoryForm = useForm<z.infer<typeof validateRepositorySchema>>({
    resolver: zodResolver(validateRepositorySchema)
  });
  const createDeploymentForm = useForm<z.infer<typeof createDeploymentSchema>>({
    resolver: zodResolver(createDeploymentSchema)
  });

  const onRepoValidationFormSubmit = async (values: z.infer<typeof validateRepositorySchema>) => {
    let { data, error } = await getRepoDetails(values.repository_url);
    if (data) {
      setRepoDetails(data);
      createDeploymentForm.setValue('name', data.name);
      createDeploymentForm.setValue('repository_url', data.svn_url);
      createDeploymentForm.setValue('branch_name', data.default_branch);
      createDeploymentForm.setValue('root_directory', '.');
      createDeploymentForm.setValue('docker_file_path', 'Dockerfile');

    } else {
      setError(error);
    }
  };

  const convertEnvToObj = (envs: EnvironmentVariableType[]): object => {
    let obj: any = {};
    envs.forEach((env) => {
      obj[env.key] = env.value;
    });
    return obj;
  };
  const onCreateDeploymentFormSubmit = async (values: z.infer<typeof createDeploymentSchema>) => {
    let envObj = convertEnvToObj(envs);
    let { data, error } = await createDeployment({
      branch_name: values.branch_name,
      docker_file_path: values.docker_file_path,
      env: envObj,
      repository_url: values.repository_url,
      root_dir: values.root_directory
    });
    if (data) {
      router.push(`/deployments/${data._id}/environments`);
    } else {
      setError(error);
    }
  };


  return (
    <div className="flex flex-col space-y-2">
      <div>
        <p className="text-2xl">Creating a new deployment</p>
      </div>

      <div>
        {!repoDetails && <div>
          {error && <ErrorAlert error={error} />}
          <Form {...validateRepositoryForm}>
            <form id="repo-url-form" onSubmit={validateRepositoryForm.handleSubmit(onRepoValidationFormSubmit)}>
              <div className="gap-4">
                <FormField control={
                  validateRepositoryForm.control
                } name="repository_url" render={({ field }) => (
                  <FormItem>
                    <FormLabel>Repository Git Url</FormLabel>
                    <FormControl>
                      <Input placeholder="Put your repository git URL here" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )} />

              </div>
              <div className="flex flex-row-reverse">
                <Button
                  disabled={loading}
                  type="submit"
                  size="sm"
                  className="my-2"
                >
                  Next
                </Button>
              </div>
            </form>
          </Form>

        </div>}
        {repoDetails && <>
          <div className="flex flex-col space-y-2 ">
            <Form {...createDeploymentForm}>
              <form id="create-deployment-form"
                    onSubmit={createDeploymentForm.handleSubmit(onCreateDeploymentFormSubmit)}>
                <div className="flex flex-row gap-2 justify-between min-w-full">
                  <FormField control={
                    createDeploymentForm.control
                  } name="name" render={({ field }) => (
                    <FormItem className="flex flex-row gap-2 w-full">
                      <FormLabel className="w-1/3 flex flex-col justify-center">Name</FormLabel>
                      <FormControl className="w-2/3">
                        <Input  {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )} />

                </div>
                <div className="flex flex-row gap-2 justify-between min-w-full">
                  <FormField control={
                    createDeploymentForm.control
                  } name="repository_url" render={({ field }) => (
                    <FormItem className="flex flex-row gap-2 w-full">
                      <FormLabel className="w-1/3 flex flex-col justify-center">Repository URL</FormLabel>
                      <FormControl className="w-2/3">
                        <Input  {...field} readOnly={true} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )} />

                </div>
                <div className="flex flex-row gap-2 justify-between min-w-full">
                  <FormField control={
                    createDeploymentForm.control
                  } name="branch_name" render={({ field }) => (
                    <FormItem className="flex flex-row gap-2 w-full">
                      <FormLabel className="w-1/3 flex flex-col justify-center">Branch Name</FormLabel>
                      <FormControl className="w-2/3">
                        <Input  {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )} />

                </div>
                <div className="flex flex-row gap-2 justify-between min-w-full">
                  <FormField control={
                    createDeploymentForm.control
                  } name="root_directory" render={({ field }) => (
                    <FormItem className="flex flex-row gap-2 w-full">
                      <FormLabel className="w-1/3 flex flex-col justify-center">Root Directory</FormLabel>
                      <FormControl className="w-2/3">
                        <Input  {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )} />

                </div>
                <div className="flex flex-row gap-2 justify-between min-w-full">
                  <FormField control={
                    createDeploymentForm.control
                  } name="docker_file_path" render={({ field }) => (
                    <FormItem className="flex flex-row gap-2 w-full">
                      <FormLabel className="w-1/3 flex flex-col justify-center">Docker File Path</FormLabel>
                      <FormControl className="w-2/3">
                        <Input  {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )} />

                </div>
                <label>Environment Variables</label>
                <div className="border 2px">
                  <EnvironmentComponent envs={envs} setEnvs={setEnvs} />
                </div>
                <div className="flex flex-row-reverse">
                  <Button
                    type="submit"
                    size="sm"
                    className="my-2"
                  >
                    Save, Redeploy
                  </Button>
                </div>

              </form>
            </Form>

          </div>
        </>}
      </div>

    </div>
  );
};
export default NewDeploymentPage;