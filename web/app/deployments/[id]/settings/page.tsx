'use client';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { NextPage } from 'next';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';

const SettingsPage: NextPage = () => {
  const { deploymentDetails } = useDeploymentContext();
  return (
    <div className="flex flex-col space-y-2">
      <h1 className="font-bold text-2xl">Settings</h1>
      <div className="border 1px solid black p-10 space-y-3">
        <h1 className="font-bold text-xl">General</h1>
        <div className="flex flex-col space-y-2 ">
          <div className="flex flex-row gap-2 justify-between">
            <div>Title</div>
            <Input className="w-3/4" value={deploymentDetails?.title} />
          </div>
        </div>
        <div className="flex flex-row-reverse">
          <Button>Save</Button>
        </div>


      </div>
      <div className="border 1px solid black p-10 space-y-3">
        <h1 className="font-bold text-xl">Build & Deploy</h1>
        <div className="flex flex-col space-y-2 ">
          <div className="flex flex-row gap-2 justify-between">
            <div>Repository</div>
            <Input className="w-3/4" value={deploymentDetails?.repository_url} />
          </div>
          <div className="flex flex-row gap-2 justify-between">
            <div>Branch</div>
            <Input className="w-3/4" value={deploymentDetails?.branch_name} />
          </div>
          <div className="flex flex-row gap-2 justify-between">
            <div>Dockerfile Path</div>
            <Input className="w-3/4" value={deploymentDetails?.docker_file_path} />
          </div>
        </div>
        <div className="flex flex-row-reverse">
          <Button>Save</Button>
        </div>


      </div>
    </div>
  );
};
export default SettingsPage;