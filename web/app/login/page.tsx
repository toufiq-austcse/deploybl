'use client';
import { NextPage } from 'next';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';

import { z } from 'zod';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAuthContext } from '@/contexts/useAuthContext';
import { useState } from 'react';
import ErrorAlert from '@/components/ui/error-alert';
import PublicRoute from '@/components/public-route';
import { useRouter } from 'next/navigation';

const formSchema = z.object({
  email: z.string().email({
    message: 'Please enter a valid email address'
  }),
  password: z.string({
    required_error: 'Password is required'
  })
});


const LoginPage: NextPage = () => {
  const { login } = useAuthContext();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema)
  });

  const onFormSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await login(values.email, values.password);
      router.push('/');
    } catch (e) {
      if (e.code === 'auth/invalid-login-credentials') {
        setError('Invalid login credentials');
      } else {
        setError(e.code);
      }
      console.log(e.code);

    }
  };

  return (
    <div className="flex m-5">
      <div className="m-auto h-1/4 w-1/4">
        <h1 className="text-4xl flex justify-center">Login</h1>
        <p className="flex justify-center">
          {`Don't have an account yet?`}
          <Link
            href={'signup'}
            className="mx-1 underline text-blue-600 hover:text-blue-800 visited:text-purple-600"
          >
            Signup
          </Link>
        </p>
        {error && <ErrorAlert error={error} />}
        <Form {...form}>
          <form id="signup-form" onSubmit={form.handleSubmit(onFormSubmit)}>
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter your email" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Enter your password"
                      type="password"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex flex-row-reverse">
              <Button
                //   disabled={loading}
                type="submit"
                size="sm"
                form="signup-form"
                className="my-2"
              >
                Login
              </Button>
            </div>
          </form>
        </Form>
      </div>
    </div>
  );

};
export default PublicRoute(LoginPage);