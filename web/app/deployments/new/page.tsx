import { NextPage } from 'next';
import { Input } from '@/components/ui/input';
import EnvironmentComponent from '@/components/ui/environment-component';

const NewDeploymentPage: NextPage = () => {
  return (
    <div className="flex flex-col space-y-2">
      <div>
        <p className="text-2xl">Creating a new deployment</p>
      </div>

      <div className="flex flex-col space-y-2 ">
        <div className="flex flex-row gap-2 justify-between">
          <div>Name</div>
          <Input className="w-3/4" />
        </div>
        <div className="flex flex-row gap-2 justify-between">
          <div>Repository</div>
          <Input className="w-3/4" />
        </div>
        <div className="flex flex-row gap-2 justify-between">
          <div>Branch</div>
          <Input className="w-3/4" />
        </div>
        <div className="flex flex-row gap-2 justify-between">
          <div>Root Directory</div>
          <Input className="w-3/4" />
        </div>
        <div className="flex flex-row gap-2 justify-between">
          <div>Dockerfile Path</div>
          <Input className="w-3/4" />
        </div>
      </div>
      <div>Environment Variables</div>
      <div className="border 2px">

        <EnvironmentComponent />
      </div>

    </div>
  );
};
export default NewDeploymentPage;