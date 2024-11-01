'use client';
import { NextPage } from 'next';
import { ColumnDef } from '@tanstack/react-table';
import AppTable from '@/components/ui/app-table';
import * as React from 'react';
import { GoCheckCircleFill } from 'react-icons/go';
// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type DeploymentEvent = {
  title: string;
  reason: string;
  latest_status: string;
  createdAt: string;
  triggeredBy: string;
};

const columns: ColumnDef<DeploymentEvent>[] = [
  {
    accessorKey: 'title',
    header: 'Title',
    cell: ({ row }) => {
      return (
        <div>
          <div className="flex flex-row gap-2">
            <GoCheckCircleFill size={16} className="mt-1" />
            <div>
              <p className="text-base font-semibold">{row.original.title}</p>
              {row.original.reason && <p className="text-sm text-gray-500">{row.original.reason}</p>}
            </div>
          </div>
        </div>
      );
    },
  },
  {
    accessorKey: 'triggeredBy',
    header: 'Triggered By',
  },
  {
    accessorKey: 'createdAt',
    header: 'Created At',
  },
];

const EventPage: NextPage = () => {
  let data: DeploymentEvent[] = [
    {
      title: 'New Deployment Started',
      reason: 'Deployment',
      latest_status: 'Deploying',
      createdAt: new Date().toDateString(),
      triggeredBy: 'User',
    },
    {
      title: 'New Deployment Started',
      reason: 'Deployment',
      latest_status: 'Deploying',
      createdAt: new Date().toDateString(),
      triggeredBy: 'User',
    },
    {
      title: 'New Deployment Started',
      reason: 'Deployment',
      latest_status: 'Deploying',
      createdAt: new Date().toDateString(),
      triggeredBy: 'User',
    },
  ];
  return (
    <div>
      <AppTable<DeploymentEvent>
        showHeader={false}
        totalPageCount={1}
        data={data}
        columns={columns}
        pageIndex={1}
        pageSize={1}
        next={() => {
          console.log('');
        }}
        prev={() => {
          console.log('');
        }}
      />
    </div>
  );
};

export default EventPage;
