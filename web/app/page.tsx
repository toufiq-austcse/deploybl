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
import { DeploymentType, PaginationType } from '@/api/http/types/deployment_type';
import DeploymentStatusBadge from '@/components/ui/deployment-status-badge';
import { IoMdAdd } from 'react-icons/io';
import { useRouter } from 'next/navigation';
import { DEPLOYMENT_STATUS } from '@/lib/constant';
import { onCopyUrlClicked } from '@/lib/utils';
import PrivateRoute from '@/components/private-route';
import { toast } from 'sonner';


const HomePage: NextPage = () => {
  const router = useRouter();
  let pageSize = Number(process.env.NEXT_PUBLIC_DEPLOYMENT_LIST_PAGE_SIZE) || 10;
  let [pageIndex, setPageIndex] = useState(1);
  let [pagination, setPagination] = useState<PaginationType>();
  let [deploymentList, setDeploymentList] = useState<DeploymentType[]>([]);
  let {
    listDeployments,
    getDeploymentLatestStatus,
    restartDeployment,
    rebuildAndDeploy,
    stopDeployment
  } = useHttpClient();
  let [loading, setLoading] = useState(true);

  const columns: ColumnDef<DeploymentType>[] = [
    {
      id: 'select',
      header: ({ table }) => (
        <Checkbox
          checked={
            (table.getIsAllPageRowsSelected() ||
              (table.getIsSomePageRowsSelected() && 'indeterminate')) as boolean
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
        const deployment = row.original;

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
              {row.getValue('latest_status') === DEPLOYMENT_STATUS.LIVE && <DropdownMenuItem
                onClick={() => onCopyUrlClicked(deployment.domain_url)}
              >
                Copy URL
              </DropdownMenuItem>}
              {row.getValue('latest_status') === DEPLOYMENT_STATUS.LIVE && <DropdownMenuItem
                onClick={() => window.open(deployment.domain_url, '_blank')}
              >
                Visit
              </DropdownMenuItem>}
              <DropdownMenuItem
                onClick={() => onRebuildAndDeployClicked(deployment._id)}
              >
                Rebuild & Deploy
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => onRestartClicked(deployment._id)}
              >
                Restart
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => onStopClicked(deployment._id)}>
                Stop
              </DropdownMenuItem>
              {row.getValue('latest_status') === DEPLOYMENT_STATUS.LIVE ? <>
              </> : null}
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <Link href={`/deployments/${deployment._id}/settings`}>
                  Settings
                </Link>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      }
    }
  ];

  useEffect(() => {
    if (deploymentList.length === 0) {
      listDeployments(pageIndex, pageSize).then(response => {
        if (response.error) {
          console.log(response.error);
        } else {
          setDeploymentList(response.data as DeploymentType[]);
          setPagination(response.pagination as PaginationType);
        }
        setLoading(false);
      });
    }

    if (deploymentList.length > 0) {
      const newIntervalId = setInterval(() => {
        updateLatestStatus(deploymentList);
      }, +(process.env.NEXT_PUBLIC_PULL_DELAY_MS as string));
      return () => clearInterval(newIntervalId);
    }


  }, [deploymentList.length]);

  const updateLatestStatus = (deployments: DeploymentType[]) => {
    let deploymentIds = deployments.map((deployment) => deployment._id);
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
    if (pageIndex === pagination?.last_page) {
      return;
    }

    listDeployments(pageIndex + 1, pageSize).then(response => {
      if (response.error) {
        console.log(response.error);
      } else {
        setPageIndex((prev) => prev + 1);
        setDeploymentList(response.data as DeploymentType[]);
        setPagination(response.pagination as PaginationType);
      }
    });

  };

  const prevFunction = () => {
    if (pageIndex === 1) {
      return;
    }
    listDeployments(pageIndex - 1, pageSize).then(response => {
      if (response.error) {
        toast.error(response.error);
        return;
      }
      setPageIndex((prev) => prev - 1);
      setDeploymentList(response.data as DeploymentType[]);
      setPagination(response.pagination as PaginationType);
    });
  };
  const onRestartClicked = async (deploymentId: string) => {
    let response = await restartDeployment(deploymentId);
    if (response.error) {
      toast.error(response.error);
      return;
    }
    toast('Deployment restarting...');
  };
  const onRebuildAndDeployClicked = async (deploymentId: string) => {
    let response = await rebuildAndDeploy(deploymentId);
    if (response.error) {
      toast.error(response.error);
      return;
    }
    toast('Deployment rebuilding and deploying...');
  };

  const onStopClicked = async (deploymentId) => {
    let response = await stopDeployment(deploymentId);
    if (response.error) {
      toast.error(response.error);
      return;
    }
    toast('Deployment stopping...');
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

        totalPageCount={pagination?.last_page as number}
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
export default PrivateRoute(HomePage);
