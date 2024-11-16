'use client';
import { NextPage } from 'next';
import { ColumnDef } from '@tanstack/react-table';
import * as React from 'react';
import { useEffect, useState } from 'react';
import { GoCheckCircleFill, GoXCircleFill } from 'react-icons/go';
import { useHttpClient } from '@/api/http/useHttpClient';
import { useDeploymentContext } from '@/contexts/useDeploymentContext';
import { DeploymentEventType, PaginationType } from '@/api/http/types/deployment_type';
import { DEPLOYMENT_EVENT_STATUS } from '@/lib/constant';
import { toast } from 'sonner';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { MoreHorizontal } from 'lucide-react';
import AppTable from '@/components/ui/app-table';
import LoggerDialog from '@/components/ui/logger-dialog';
import { formatDateTime } from '@/lib/utils';

const EventPage: NextPage = () => {
  const { deploymentDetails } = useDeploymentContext();
  const { listDeploymentEvents } = useHttpClient();
  const [loading, setLoading] = React.useState(false);
  let pageSize = Number(process.env.NEXT_PUBLIC_DEPLOYMENT_EVENT_LIST_PAGE_SIZE) || 10;
  let [pageIndex, setPageIndex] = useState(1);
  let [pagination, setPagination] = useState<PaginationType>();
  const [columns, setColumns] = useState<ColumnDef<DeploymentEventType>[]>([]);
  const [isOpenLoggerDialog, setIsOpenLoggerDialog] = useState(false);
  const [loggingEvent, setLoggingEvent] = useState<DeploymentEventType | null>(null);

  const [deploymentEvents, setDeploymentEvents] = React.useState<DeploymentEventType[]>([]);

  const getDeploymentEventIcon = (status: string) => {
    switch (status) {
      case DEPLOYMENT_EVENT_STATUS.SUCCESS:
        return <GoCheckCircleFill className="text-green-800/80 mt-1" size={16} />;
      case DEPLOYMENT_EVENT_STATUS.FAILED:
        return <GoXCircleFill className="text-destructive mt-1" size={16} />;
      default:
        return (
          <div className="mt-1">
            <svg
              aria-hidden="true"
              className="w-4 h-4 text-gray-200 animate-spin dark:text-gray-600 fill-blue-600"
              viewBox="0 0 100 101"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                fill="currentColor"
              />
              <path
                d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                fill="currentFill"
              />
            </svg>
            <span className="sr-only">Loading...</span>
          </div>
        );
    }
  };

  useEffect(() => {
    if (!deploymentDetails) {
      return;
    }
    if (deploymentEvents.length == 0) {
      setLoading(true);
    }
    listDeploymentEvents(deploymentDetails._id, pageIndex, pageSize)
      .then((response) => {
        if (!response.isSuccessful && response.code !== 401) {
          toast.error(response.error);
          return;
        }
        setDeploymentEvents(response.data);
        setPagination(response.pagination);
        initColumns();
      })
      .finally(() => {
        setLoading(false);
      });
  }, [deploymentDetails.latest_status]);

  const initColumns = () => {
    setColumns([
      {
        accessorKey: 'title',
        header: 'Title',
        cell: ({ row }) => {
          let latestStatus = row.original.status;
          return (
            <div className="flex flex-row gap-2">
              {getDeploymentEventIcon(latestStatus)}
              <div>
                <p className="text-base font-semibold">{row.original.title}</p>
                {row.original.reason && <p className="text-sm text-gray-500">{row.original.reason}</p>}
              </div>
            </div>
          );
        },
      },
      {
        accessorKey: 'triggered_by',
        header: 'Triggered By',
      },
      {
        accessorKey: 'created_at',
        header: 'Created At',
        cell: ({ row }) => {
          return <div>{formatDateTime(row.original.created_at)}</div>;
        },
      },
      {
        id: 'actions',
        enableHiding: false,
        cell: ({ row }) => {
          const event = row.original;

          return (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem
                  onClick={() => {
                    setLoggingEvent(event);
                    setIsOpenLoggerDialog(true);
                  }}
                >
                  View logs
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          );
        },
      },
    ]);
  };
  const nextFunction = () => {
    if (pageIndex === pagination?.last_page) {
      return;
    }

    listDeploymentEvents(deploymentDetails._id, pageIndex + 1, pageSize).then((response) => {
      if (!response.isSuccessful && response.code !== 401) {
        toast.error(response.error);
        return;
      }
      setPageIndex((prev) => prev + 1);
      setDeploymentEvents(response.data);
      setPagination(response.pagination as PaginationType);
    });
  };

  const prevFunction = () => {
    if (pageIndex === 1) {
      return;
    }
    listDeploymentEvents(deploymentDetails._id, pageIndex - 1, pageSize).then((response) => {
      if (!response.isSuccessful && response.code !== 401) {
        toast.error(response.error);
        return;
      }
      setPageIndex((prev) => prev - 1);
      setDeploymentEvents(response.data);
      setPagination(response.pagination as PaginationType);
    });
  };

  return (
    <div>
      {loading ? (
        <p>Loading...</p>
      ) : (
        <div>
          <LoggerDialog open={isOpenLoggerDialog} setOpen={setIsOpenLoggerDialog} loggingEvent={loggingEvent} />
          <AppTable<DeploymentEventType>
            showHeader={false}
            showCaption={false}
            totalPageCount={pagination?.last_page}
            data={deploymentEvents}
            columns={columns}
            pageIndex={pageIndex}
            pageSize={pageSize}
            next={nextFunction}
            prev={prevFunction}
          />
        </div>
      )}
    </div>
  );
};

export default EventPage;
