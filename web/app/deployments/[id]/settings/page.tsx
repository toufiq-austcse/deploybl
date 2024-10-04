'use client';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { NextPage } from 'next';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
import * as React from 'react';
import { useEffect, useState } from 'react';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form';
import ErrorAlert from '@/components/ui/error-alert';
import { toast } from 'sonner';

const generalUpdateSchema = z.object({
  title: z.string({
    required_error: 'Title is required'
  }).min(1, {
    message: 'title is required'
  })
});
const buildAndDeploySchema = z.object({
  branch_name: z.string({
    required_error: 'Branch name is required'
  }),
  root_directory: z.string().optional().default(''),
  docker_file_path: z.string({
    required_error: 'Docker file path is required'
  }).min(1, {
    message: 'Docker file path is required'
  })
});
const SettingsPage: NextPage = () => {
  const { deploymentDetails, updateDeploymentDetails } = useDeploymentContext();
  const [latestDeploymentDetails] = useState(deploymentDetails);
  const [loading, setLoading] = useState(false);
  const [generalError, serGeneralError] = useState<string | null>(null);
  const [buildAndDeployError, setBuildAndDeployError] = useState<string | null>(null);

  const validateGeneralUpdateForm = useForm<z.infer<typeof generalUpdateSchema>>({
    resolver: zodResolver(generalUpdateSchema),
    mode: 'onChange'
  });
  const validateBuildAndDeployForm = useForm<z.infer<typeof buildAndDeploySchema>>({
    resolver: zodResolver(buildAndDeploySchema),
    mode: 'onChange'
  });

  useEffect(() => {
    validateGeneralUpdateForm.setValue('title', deploymentDetails?.title);
    validateBuildAndDeployForm.setValue('root_directory', deploymentDetails?.root_directory == null ? '' : deploymentDetails?.root_directory);
    validateBuildAndDeployForm.setValue('branch_name', deploymentDetails?.branch_name);
    validateBuildAndDeployForm.setValue('docker_file_path', deploymentDetails?.docker_file_path);
  }, [latestDeploymentDetails]);

  const onGeneralUpdateFormSubmit = async (values: z.infer<typeof generalUpdateSchema>) => {
    setLoading(true);
    serGeneralError(null);
    updateDeploymentDetails(deploymentDetails?._id, { title: values.title }).then(response => {
      if (response.error) {
        serGeneralError(response.error);
      } else {
        toast('Successfully updated');
      }
    }).finally(() => {
      setLoading(false);
    });
  };
  const onBuildAndDeployFormSubmit = async (values: z.infer<typeof buildAndDeploySchema>) => {
    setLoading(true);
    setBuildAndDeployError(null);
    updateDeploymentDetails(deploymentDetails?._id, {
      branch_name: values.branch_name,
      root_dir: values.root_directory === '' ? null : values.root_directory,
      docker_file_path: values.docker_file_path
    }).then(response => {
      if (response.error) {
        setBuildAndDeployError(response.error);
      } else {
        toast('Successfully updated');
      }
    }).finally(() => {
      setLoading(false);
    });
  };

  return (
    <div className="flex flex-col space-y-2">
      <h1 className="font-bold text-2xl">Settings</h1>
      <div className="border 1px solid black p-10 space-y-3">
        <h1 className="font-bold text-xl">General</h1>
        {generalError && <ErrorAlert error={generalError} />}
        <Form {...validateGeneralUpdateForm}>
          <form id="repo-url-form" onSubmit={validateGeneralUpdateForm.handleSubmit(onGeneralUpdateFormSubmit)}>
            <div className="gap-4">
              <FormField control={
                validateGeneralUpdateForm.control
              } name="title" render={({ field }) => (
                <FormItem>
                  <div className="flex flex-col space-y-2">
                    <div className="flex flex-row gap-2 justify-between">
                      <div className="w-1/3">Title</div>
                      <div className="flex flex-col w-full">
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </div>
                    </div>
                  </div>

                </FormItem>
              )} />

            </div>
            <div className="flex flex-row-reverse">
              <Button
                disabled={!validateGeneralUpdateForm.formState.isDirty || loading}
                type="submit"
                size="sm"
                className="my-2"
              >
                Save
              </Button>
            </div>
          </form>
        </Form>
      </div>
      <div className="border 1px solid black p-10 space-y-3">
        <h1 className="font-bold text-xl">Build & Deploy</h1>
        {buildAndDeployError && <ErrorAlert error={buildAndDeployError} />}
        <Form {...validateBuildAndDeployForm}>
          <form id="create-deployment-form"
                onSubmit={validateBuildAndDeployForm.handleSubmit(onBuildAndDeployFormSubmit)}>
            <div className="my-2">
              <FormField control={
                validateBuildAndDeployForm.control
              } name="branch_name" render={({ field }) => (
                <FormItem>
                  <div className="flex flex-col space-y-2">
                    <div className="flex flex-row gap-2 justify-between">
                      <div className="w-1/3">Branch</div>
                      <div className="flex flex-col w-full">
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </div>
                    </div>
                  </div>

                </FormItem>

              )} />
            </div>
            <div className="my-2">
              <FormField control={
                validateBuildAndDeployForm.control
              } name="root_directory" render={({ field }) => (
                <FormItem>
                  <div className="flex flex-col space-y-2">
                    <div className="flex flex-row gap-2 justify-between">
                      <div className="w-1/3">Root Directory</div>
                      <div className="flex flex-col w-full">
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </div>
                    </div>
                  </div>

                </FormItem>

              )} />

            </div>
            <div>
              <FormField control={
                validateBuildAndDeployForm.control
              } name="docker_file_path" render={({ field }) => (
                <FormItem>
                  <div className="flex flex-col space-y-2">
                    <div className="flex flex-row gap-2 justify-between">
                      <div className="w-1/3">Dockerfile Path</div>
                      <div className="flex flex-col w-full">
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </div>
                    </div>
                  </div>

                </FormItem>

              )} />
            </div>
            <div className="flex flex-row-reverse">
              <Button
                disabled={loading || !validateBuildAndDeployForm.formState.isDirty}
                type="submit"
                size="sm"
                className="my-2"
              >
                Save, Rebuild & Redeploy
              </Button>
            </div>
          </form>
        </Form>

      </div>
    </div>
  );
};
export default SettingsPage;