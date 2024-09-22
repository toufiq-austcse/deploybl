import * as React from 'react';

const Loading = ({ className='bg-black' }: { className: string }) => {
  return (
    <div className="flex flex-row">
      <div className={`h-1.5 w-1.5 ${className} bg-white rounded-full animate-bounce [animation-delay:-0.3s]`}></div>
      <div className={`h-1.5 w-1.5 ${className} rounded-full animate-bounce [animation-delay:-0.15s]`}></div>
      <div className={`h-1.5 w-1.5 ${className} rounded-full animate-bounce`}></div>
    </div>
  );
};
export default Loading;