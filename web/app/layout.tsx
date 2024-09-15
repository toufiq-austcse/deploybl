'use client';
import '@/styles/globals.css';
import Navbar from '@/components/ui/navbar';
import * as React from 'react';

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <html>
    <body>
    <div className={'min-h-screen flex flex-col'}>
      <Navbar />
      <div className="m-4 max-w-full px-5 sm:px-6 lg:px-20">
        {children}
      </div>
    </div>
    </body>
    </html>
  );
};
export default RootLayout;