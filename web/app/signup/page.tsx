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
import { useEffect, useState } from 'react';
import ErrorAlert from '@/components/ui/error-alert';
import PublicRoute from '@/components/public-route';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { FaEnvelope, FaGithub, FaLock, FaUser } from 'react-icons/fa';
import { FcGoogle } from 'react-icons/fc';

const formSchema = z.object({
  name: z.string({
    required_error: 'Name is required'
  }),
  email: z.string().email({
    message: 'Please enter a valid email address'
  }),
  password: z.string({
    required_error: 'Password is required'
  })
});

const SignUpPage: NextPage = () => {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const { signup, loginWithGithub, loginWithGoogle } = useAuthContext();
  const [error, setError] = useState<string | null>(null);
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema)
  });

  // Add subtle animation effect when component mounts
  useEffect(() => {
    // This is where you could add animations if needed
    // For example, fade-in effects or loading states
  }, []);

  const onFormSubmit = async (values: z.infer<typeof formSchema>) => {
    setLoading(true);
    try {
      await signup(values.email, values.password, values.name);
      router.push('/');
    } catch (err) {
      if (err.code === 'auth/email-already-in-use') {
        setError('Email already in use');
      } else if (err.code === 'auth/weak-password') {
        setError('Weak password');
      } else {
        setError(err.code);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex justify-center">
      <div className="max-w-lg w-full space-y-4 px-4 py-6 sm:px-6 lg:px-8 flex flex-col items-center">
        {/* Logo/Branding */}
        <div className="text-center mb-2">
          <h2 className="mt-2 text-center text-4xl font-extrabold text-gray-900">Sign Up</h2>
          <p className="mt-2 text-center text-base text-gray-600">
            Already have an account?{' '}
            <Link href={'login'}
                  className="font-medium text-blue-600 hover:text-blue-500 transition-colors duration-200">
              Login
            </Link>
          </p>
        </div>

        <Card
          className="shadow-2xl border border-gray-100 overflow-hidden mx-auto w-full rounded-lg transform transition-all duration-300 hover:shadow-3xl">
          {error && <div className="px-6 pt-6"><ErrorAlert error={error} /></div>}

          <CardContent className="pt-6 px-8 pb-6">
            {/* GitHub Login Button */}
            <div className="mb-3">
              <Button
                onClick={loginWithGithub}
                className="w-full flex items-center justify-center gap-2 bg-gray-800 hover:bg-gray-900 transition-colors duration-200 py-3 rounded-md"
              >
                <FaGithub className="text-lg" />
                <span className="font-medium">Continue with GitHub</span>
              </Button>
            </div>

            {/* Google Login Button */}
            <div className="mb-4">
              <Button
                onClick={loginWithGoogle}
                className="w-full flex items-center justify-center gap-2 bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 transition-colors duration-200 py-3 rounded-md"
              >
                <FcGoogle className="text-lg" />
                <span className="font-medium">Continue with Google</span>
              </Button>
            </div>

            {/* Divider */}
            <div className="relative mb-4">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-200"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-white text-gray-500">Or continue with email</span>
              </div>
            </div>

            {/* Signup Form */}
            <Form {...form}>
              <form id="signup-form" onSubmit={form.handleSubmit(onFormSubmit)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-gray-700 text-sm font-medium">Name</FormLabel>
                      <div className="relative">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                          <FaUser className="h-5 w-5 text-gray-400" />
                        </div>
                        <FormControl>
                          <Input
                            placeholder="Enter your name"
                            className="pl-10 py-3 bg-gray-50 border-gray-300 focus:ring-blue-500 focus:border-blue-500 rounded-md text-base"
                            {...field}
                          />
                        </FormControl>
                      </div>
                      <FormMessage className="text-red-500" />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-gray-700 text-sm font-medium">Email</FormLabel>
                      <div className="relative">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                          <FaEnvelope className="h-5 w-5 text-gray-400" />
                        </div>
                        <FormControl>
                          <Input
                            placeholder="Enter your email"
                            className="pl-10 py-3 bg-gray-50 border-gray-300 focus:ring-blue-500 focus:border-blue-500 rounded-md text-base"
                            {...field}
                          />
                        </FormControl>
                      </div>
                      <FormMessage className="text-red-500" />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-gray-700 text-sm font-medium">Password</FormLabel>
                      <div className="relative">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                          <FaLock className="h-5 w-5 text-gray-400" />
                        </div>
                        <FormControl>
                          <Input
                            placeholder="Enter your password"
                            type="password"
                            className="pl-10 py-3 bg-gray-50 border-gray-300 focus:ring-blue-500 focus:border-blue-500 rounded-md text-base"
                            {...field}
                          />
                        </FormControl>
                      </div>
                      <FormMessage className="text-red-500" />
                    </FormItem>
                  )}
                />

                <Button
                  disabled={loading}
                  type="submit"
                  form="signup-form"
                  className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 rounded-md transition-colors duration-200 mt-4 text-base"
                >
                  {loading ? 'Signing up...' : 'Sign up'}
                </Button>
              </form>
            </Form>
          </CardContent>

          <CardFooter className="bg-gray-50 px-8 py-4 flex justify-center border-t border-gray-100">
            <p className="text-sm text-gray-500">
              By signing up, you agree to our <Link href="#"
                                                    className="text-blue-600 hover:text-blue-500 font-medium">Terms</Link> and <Link
              href="#" className="text-blue-600 hover:text-blue-500 font-medium">Privacy</Link>
            </p>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
};
export default PublicRoute(SignUpPage);
