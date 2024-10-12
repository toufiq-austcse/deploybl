'use client';
import { MdAdd, MdDelete } from 'react-icons/md';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { IoMdAdd } from 'react-icons/io';

export interface EnvironmentVariableType {
  key: string;
  value: string;
}

const EnvironmentComponent = ({ envs = [], setEnvs }: { envs: EnvironmentVariableType[]; setEnvs: Function }) => {
  const handleKeyChange = (e: any, index: number) => {
    setEnvs((prev: any) => {
      prev[index].key = e.target.value;
      return [...prev];
    });
  };

  const handleValueChange = (e: any, index: number) => {
    setEnvs((prev: any) => {
      prev[index].value = e.target.value;
      return [...prev];
    });
  };

  return (
    <div className="m-2">
      {envs.length > 0 && (
        <div className="flex flex-row m-2">
          <h1 className="w-1/2">Key</h1>
          <h1 className="w-1/2">Value</h1>
        </div>
      )}
      {envs.map((value, index, array) => {
        return (
          <div key={index} className="flex flex-row m-2 gap-4">
            <div className="w-1/2">
              <Input
                placeholder="Name of the varibale"
                value={envs[index].key}
                onChange={(e) => handleKeyChange(e, index)}
              />
            </div>
            <div className="w-1/2">
              <Input placeholder="value" value={envs[index].value} onChange={(e) => handleValueChange(e, index)} />
            </div>

            <div className="flex flex-col justify-center">
              <MdDelete
                onClick={() => {
                  envs.splice(index, 1);
                  console.log('newEnvs ', envs);
                  setEnvs(() => [...envs]);
                }}
              />
            </div>
          </div>
        );
      })}
      <div className="flex flex-row m-2">
        <Button
          onClick={(event) => {
            event.preventDefault();
            setEnvs((prev: any) => [...prev, {}]);
          }}
          size="sm"
          variant="outline"
        >
          <div className="flex flex-row justify-between gap-2">
            <div className="flex flex-col justify-center">
              <IoMdAdd />
            </div>
            <div>Add environment variable</div>
          </div>
        </Button>
      </div>
    </div>
  );
};
export default EnvironmentComponent;
