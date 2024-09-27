'use client';

import * as React from 'react';
import { useEffect, useState } from 'react';
import { ColumnDef } from '@tanstack/react-table';
import { MoreHorizontal } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import Link from 'next/link';
import AppTable from '@/components/ui/app-table';
import { NextPage } from 'next';
import { useHttpClient } from '@/api/http/useHttpClient';
import { DeploymentType } from '@/api/http/types/deployment_type';
import DeploymentStatusBadge from '@/components/ui/deployment-status-badge';
import { IoMdAdd } from 'react-icons/io';
import { useRouter } from 'next/navigation';
import { DEPLOYMENT_STATUS } from '@/lib/constant';


const columns: ColumnDef<DeploymentType>[] = [
  {
    id: 'select',
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && 'indeterminate')
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false
  },

  {
    accessorKey: 'title',
    header: 'Title',
    cell: ({ row }) => {
      return (
        <Link className="flex space-x-4" href={`/deployments/${row.original._id}/settings`}>
          <div className="lowercase font-medium">{row.getValue('title')}</div>
        </Link>
      );
    }
  },
  {
    accessorKey: 'latest_status',
    header: 'Status',
    cell: ({ row }) => {
      let status: string = row.getValue('latest_status');
      return <DeploymentStatusBadge status={status} />;
    }
  },
  {
    accessorKey: 'repository_provider',
    header: 'Provider',
    cell: ({ row }) => {
      return <div>{row.getValue('repository_provider')}</div>;
    }
  },
  {
    accessorKey: 'last_deployed_at',
    header: 'Last Deployed At',
    cell: ({ row }) => {
      let lastDeployedAt = row.getValue('last_deployed_at');
      if (!lastDeployedAt) {
        return <div>Not Deployed Yet</div>;
      }
      const date = new Date(row.getValue('last_deployed_at')).toDateString();

      return <div>{date}</div>;
    }
  },
  {
    id: 'actions',
    enableHiding: false,
    cell: ({ row }) => {
      const video = row.original;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(video._id)}
            >
              Copy Url
            </DropdownMenuItem>
            <DropdownMenuSeparator />
          </DropdownMenuContent>
        </DropdownMenu>
      );
    }
  }
];

const HomePage: NextPage = () => {
  const router = useRouter();
  let pageSize = Number(process.env.NEXT_PUBLIC_VIDEO_LIST_PAGE_SIZE) || 4;
  let [pageIndex, setPageIndex] = useState(0);
  let [deploymentList, setDeploymentList] = useState<DeploymentType[]>([]);
  let { listDeployments, getDeploymentLatestStatus } = useHttpClient();
  let [isInitialLoadingDone, setIsInitialLoadingDone] = useState(false);
  let [loading, setLoading] = useState(true);

  useEffect(() => {
    if (deploymentList.length === 0 && !isInitialLoadingDone) {
      listDeployments(0, pageSize).then(response => {
        if (response.error) {
          console.log(response.error);
        } else {
          setDeploymentList(response.data as DeploymentType[]);
        }
        setIsInitialLoadingDone(true);
        setLoading(false);
      });
    }

    if (deploymentList.length > 0) {
      const newIntervalId = setInterval(() => {
        updateLatestStatus(deploymentList);
      }, +(process.env.NEXT_PUBLIC_PULL_DELAY_MS as string));
      return () => clearInterval(newIntervalId);
    }


  }, [deploymentList]);

  const updateLatestStatus = (deployments: DeploymentType[]) => {
    let deploymentIds = deployments.filter(deployment => deployment.latest_status !== DEPLOYMENT_STATUS.LIVE && deployment.latest_status !== DEPLOYMENT_STATUS.FAILED).map((deployment) => deployment._id);
    if (deploymentIds.length === 0) {
      return;
    }
    getDeploymentLatestStatus(deploymentIds).then(response => {
      if (response.error) {
        console.log(response.error);
      } else {
        console.log(response);
        // @ts-ignore
        if (response.data?.length > 0) {
          console.log('setting...');
          setDeploymentList((deployments) => {
            return deployments.map((deployment) => {
              // @ts-ignore
              let latestStatus = response.data.find((status) => status._id === deployment._id);
              if (latestStatus) {
                return {
                  ...deployment,
                  latest_status: latestStatus.latest_status,
                  last_deployed_at: latestStatus.last_deployed_at
                };
              }
              return deployment;
            });
          });
        }
      }
    });


  };

  const nextFunction = () => {
    console.log('next');
    // fetchMore({
    //   variables: {
    //     first: pageSize,
    //     after: data.ListAsset.page_info.next_cursor,
    //   },
    //   updateQuery: (prev, { fetchMoreResult }) => {
    //     setPageIndex((prev) => prev + 1);
    //     if (!fetchMoreResult) {
    //       return prev;
    //     }
    //     return fetchMoreResult;
    //   },
    // });
  };

  const prevFunction = () => {
    console.log('prev');
    // fetchMore({
    //   variables: {
    //     first: pageSize,
    //     before: data.ListAsset.page_info.prev_cursor,
    //   },
    //   updateQuery: (prev, { fetchMoreResult }) => {
    //     setPageIndex((prev) => prev - 1);
    //     if (!fetchMoreResult) {
    //       return prev;
    //     }
    //     return fetchMoreResult;
    //   },
    // });
  };

  return (

    <div className="space-y-2">
      <div className="flex flex-row-reverse gap-2">
        <div className="flex flex-row">
          <Button variant="outline" onClick={() => {
            router.push('/deployments/new');
          }}>
            <div className="flex flex-row justify-between gap-2">
              <div className="flex flex-col justify-center">
                <IoMdAdd />
              </div>
              <div>Create New</div>
            </div>
          </Button>
        </div>
      </div>
      {loading ? <div className="justify-end">Loading...</div> : <AppTable<DeploymentType>

        totalPageCount={0}
        data={deploymentList}
        columns={columns}
        pageIndex={pageIndex}
        pageSize={pageSize}
        next={nextFunction}
        prev={prevFunction}
      />}

    </div>
  );
};
export default HomePage;
