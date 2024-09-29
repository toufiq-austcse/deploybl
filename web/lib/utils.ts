import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { EnvironmentVariableType } from '@/components/ui/environment-component';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function bytesToMegaBytes(bytes: number) {
  return (bytes / 1024 / 1024).toFixed(2);
}

export function secondsToHHMMSS(seconds: number) {
  return new Date(seconds * 1000).toISOString().substr(11, 8);
}

export const convertEnvToObj = (envs: EnvironmentVariableType[]): object => {
  let obj: any = {};
  envs.forEach((env) => {
    if (env.key && env.value) {
      obj[env.key] = env.value;
    }

  });
  return obj;
};
export const convertObjToEnv = (envs: object): EnvironmentVariableType[] => {
  let arr: EnvironmentVariableType[] = [];
  for (let key in envs) {
    arr.push({ key, value: envs[key] });
  }
  return arr;
};