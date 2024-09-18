'use client';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { MdDelete } from 'react-icons/md';
import { Button } from '@/components/ui/button';

const EnvironmentComponent = () => {
  const [envs, setEnvs] = useState([{}]);
  const form = useForm();

  return (
    <div>
      <Form {...form}>
        <form id="login-form">
          <div className="flex flex-row gap-2 justify-center m-2">
            <h1 className="w-2/5">Key</h1>
            <h1 className="w-2/5">Value</h1>
          </div>

          {envs.map((value, index, array) => {
            return (
              <div key={index} className="flex flex-row gap-2 justify-center m-2">
                <div className="w-2/5">
                  <FormField
                    control={form.control}
                    name="key"
                    render={({ field }) => (
                      <FormItem>
                        <FormControl>
                          <Input placeholder="Name of variable" {...field} />
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
                          <Input placeholder="Name of variable" {...field} />
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

              </div>
            );
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
export default EnvironmentComponent;