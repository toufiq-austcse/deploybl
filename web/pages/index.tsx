'use client';

import * as React from 'react';
import { useEffect } from 'react';
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
import { badgeVariants } from '@/components/ui/badge';
import Link from 'next/link';
import AppTable from '@/components/ui/app-table';
import { DEPLOYMENT_STATUS } from '@/lib/constant';
import { NextPage } from 'next';
import { useHttpClient } from '@/api/http/useHttpClient';
import { DeploymentType } from '@/api/http/types/deployment_type';

export type Deployment = {
  _id: string;
  title: string;
  last_deployed_at: string;
  repository_provider: string
  latest_status: 'pending' | 'processing' | 'success' | 'failed';
  created_at: Date;
  updated_at: Date;
};

export const columns: ColumnDef<DeploymentType>[] = [
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
        <Link className="flex space-x-4" href={`/deployments/${row.original._id}`}>
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
      if (status === DEPLOYMENT_STATUS.FAILED) {

        return (
          <div
            className={`${badgeVariants({ variant: 'destructive' })} capitalize`}
          >
            {status}
          </div>
        );
      } else if (status === DEPLOYMENT_STATUS.PULLING) {
        return (
          <div
            className={`${badgeVariants({ variant: 'default' })} capitalize`}
          >
            {status}
          </div>
        );
      } else if (status === DEPLOYMENT_STATUS.BUILDING) {
        return (
          <div
            className={`${badgeVariants({ variant: 'outline' })} capitalize`}
          >
            {status}
          </div>
        );
      }

      return (
        <div
          className={`${badgeVariants({ variant: 'secondary' })} capitalize`}
        >
          {status}
        </div>
      );
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
  let pageSize = Number(process.env.NEXT_PUBLIC_VIDEO_LIST_PAGE_SIZE) || 4;
  let [pageIndex, setPageIndex] = React.useState(0);
  let [deploymentList, setDeploymentList] = React.useState<DeploymentType[]>([]);
  let { listDeployments } = useHttpClient();

  useEffect(() => {
    listDeployments(1, 0).then(data => {
      setDeploymentList(data);
    }).catch(err => console.log('error in list ', err));
  }, []);

  // let { data, loading, error, fetchMore, refetch } = useQuery(LIST_ASSETS, {
  //   variables: {
  //     first: pageSize,
  //     after: null,
  //   },
  // });

  //console.log("data ", data);

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
    <div>
      <div className="flex">

      </div>

      <AppTable<DeploymentType>
        totalPageCount={0}
        data={deploymentList}
        columns={columns}
        pageIndex={pageIndex}
        pageSize={pageSize}
        next={nextFunction}
        prev={prevFunction}
      />
    </div>
  );
};
export default HomePage;
