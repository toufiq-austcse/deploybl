import { NextPage } from 'next';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

const EnvironmentPage: NextPage = () => {
  return (
    <div className="flex flex-col space-y-2">
      <h1 className="font-bold text-2xl">Settings</h1>
      <div className="border 1px solid black p-10 space-y-3">
        <h1 className="font-bold text-xl">General</h1>
        <div className="flex flex-col space-y-2 ">
          <div className="flex flex-row gap-2 justify-between">
            <div>Name</div>
            <Input className="w-3/4" />
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
        <div className="flex flex-row-reverse">
          <Button>Save</Button>
        </div>


      </div>
    </div>
  );
};
export default EnvironmentPage;