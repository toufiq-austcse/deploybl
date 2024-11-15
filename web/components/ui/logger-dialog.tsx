'use client';
import { DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Dialog } from '@radix-ui/react-dialog';
import * as React from 'react';
import SyntaxHighlighter from 'react-syntax-highlighter';
import { dracula } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { DeploymentEventType } from '@/api/http/types/deployment_type';
import { formatDateTime } from '@/lib/utils';
import { useHttpClient } from '@/api/http/useHttpClient';

const LoggerDialog = ({
  open,
  setOpen,
  loggingEvent,
}: {
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  loggingEvent: DeploymentEventType;
}) => {
  const { getLogContents } = useHttpClient();
  const [logContent, setLogContent] = React.useState<string>('');
  React.useEffect(() => {
    if (loggingEvent) {
      getLogContents(loggingEvent.event_log_file_url).then((response) => {
        if (response.isSuccessful) {
          setLogContent(response.data);
        }
      });
    }
  }, [loggingEvent]);
  return (
    loggingEvent && (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-5xl w-full">
          <DialogHeader>
            <DialogTitle>{loggingEvent.title} log</DialogTitle>
            <DialogDescription>
              {loggingEvent.reason} - {formatDateTime(loggingEvent.created_at)}
            </DialogDescription>
          </DialogHeader>
          <div className="rounded-md shadow-md overflow-auto max-h-96">
            <SyntaxHighlighter
              language="javascript"
              style={dracula}
              customStyle={{
                fontSize: '16px',
              }}
            >
              {logContent}
            </SyntaxHighlighter>
          </div>
        </DialogContent>
      </Dialog>
    )
  );
};

export default LoggerDialog;
