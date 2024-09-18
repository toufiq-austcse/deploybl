'use client';
import * as React from 'react';
import { useForm } from 'react-hook-form';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRouter } from 'next/navigation';
import { IoMdAdd } from 'react-icons/io';

const CreateNewModal = () => {
  const router = useRouter();
  const form = useForm();
  return (
    <Dialog>
      <DialogTrigger className="btn" asChild>
        <Button variant="outline"> <IoMdAdd /> Create New</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>New Deployment</DialogTitle>
        </DialogHeader>
        <div>
          <div className=" gap-4">
            <Input
              id="name"
              placeholder="Put your repository url here"
            />
          </div>
        </div>
        <DialogFooter>
          <Button onClick={() => {
            router.push('/deployments/new');
          }} type="submit">Next</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default CreateNewModal;