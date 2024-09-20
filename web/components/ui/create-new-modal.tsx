'use client';
import * as React from 'react';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRouter } from 'next/navigation';
import { IoMdAdd } from 'react-icons/io';
import { useState } from 'react';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { useHttpClient } from '@/api/http/useHttpClient';
import ErrorAlert from '@/components/ui/error-alert';

const formSchema = z.object({
  repository_url: z
    .string()
    .url({
      message: 'Please enter valid url'
    })
});

const CreateNewModal = () => {
  const router = useRouter();
  const [error, setError] = React.useState<string | null>(null);
  const { getRepoDetails, loading } = useHttpClient();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema)
  });

  const onFormSubmit = async (values: z.infer<typeof formSchema>) => {
    console.log('values ', values);
    let { data, error } = await getRepoDetails(values.repository_url);
    if (data) {
      console.log('data ', data);
      router.push('/deployments/new');
    } else {
      setError(error);
    }
    // setLoading(true);
    // let { data, error } = await userLogin(values.email, values.password);
    // setLoading(false);
    // if (data) {
    //   localStorage.setItem("token", data.token.access_token);
    //   location.reload();
    // } else {
    //   setError(error);
    // }
  };
  return (

    <Dialog>
      <DialogTrigger className="btn" asChild>
        <Button variant="outline"> <IoMdAdd /> Create New</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>New Deployment</DialogTitle>
        </DialogHeader>
        {error && <ErrorAlert error={error} />}
        <Form {...form}>
          <form id="repo-url-form" onSubmit={form.handleSubmit(onFormSubmit)}>
            <div className="gap-4">
              <FormField control={
                form.control
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
                Login
              </Button>
            </div>
          </form>

        </Form>
      </DialogContent>
    </Dialog>


  );
};

export default CreateNewModal;