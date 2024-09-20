'use client';
import { MdDelete } from 'react-icons/md';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export interface EnvironmentVariableType {
  key: string;
  value: string;
}

const EnvironmentComponent = ({ envs = [], setEnvs }: { envs: EnvironmentVariableType[], setEnvs: Function }) => {
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
    <div>
      <div className="flex flex-row gap-2 justify-center m-2">
        <h1 className="w-2/5">Key</h1>
        <h1 className="w-2/5">Value</h1>
      </div>

      {envs.map((value, index, array) => {
        return (
          <div key={index} className="flex flex-row gap-2 justify-center m-2">
            <div className="w-2/5">
              <Input placeholder="Key of the varibale" value={envs[index].key}
                     onChange={(e) => handleKeyChange(e, index)} />
            </div>
            <div className="w-2/5">
              <Input placeholder="Value of the variable" value={envs[index].value}
                     onChange={(e) => handleValueChange(e, index)} />
            </div>

            <div className="flex flex-col justify-center">
              <MdDelete onClick={() => {
                envs.splice(index, 1);
                console.log('newEnvs ', envs);
                setEnvs(() => [...envs]);
              }} />
            </div>

          </div>
        );
      })}
      <div className="flex flex-row-reverse gap-2">

        <Button
          onClick={() => {
            setEnvs((prev: any) => [...prev, {}]);
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