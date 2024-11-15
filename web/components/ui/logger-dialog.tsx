'use client';
import { DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Dialog } from '@radix-ui/react-dialog';
import * as React from 'react';
import SyntaxHighlighter from 'react-syntax-highlighter';
import { dracula } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { DeploymentEventType } from '@/api/http/types/deployment_type';
import { formatDateTime } from '@/lib/utils';

const LoggerDialog = ({
  open,
  setOpen,
  loggingEvent,
}: {
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  loggingEvent: DeploymentEventType;
}) => {
  const logContent = `
Cloning repository from https://github.com/toufiq-austcse/test-deployitd.git branch master
cloned successfully
building dockerfile
#0 building with "desktop-linux" instance using docker driver

#1 [internal] load .dockerignore
#1 transferring context: 2B done
#1 DONE 0.0s

#2 [internal] load build definition from Dockerfile
#2 transferring dockerfile: 859B done
#2 DONE 0.0s

#3 [internal] load metadata for docker.io/library/node:16
#3 DONE 1.2s

#4 [internal] load metadata for docker.io/library/node:16-alpine
#4 DONE 1.2s

#5 [base 1/6] FROM docker.io/library/node:16@sha256:f77a1aef2da8d83e45ec990f45df50f1a286c5fe8bbfb8c6e4246c6389705c0b
#5 DONE 0.0s

#6 [stage-2 1/4] FROM docker.io/library/node:16-alpine@sha256:a1f9d027912b58a7c75be7716c97cfbc6d3099f3a97ed84aa490be9dee20e787
#6 DONE 0.0s

#7 [internal] load build context
#7 transferring context: 18.99kB done
#7 DONE 0.0s

#8 [dev 1/4] COPY nest-cli.json   tsconfig.*    ormconfig.ts   ./
#8 CACHED

#9 [dev 4/4] RUN yarn build
#9 CACHED

#10 [base 2/6] WORKDIR /app
#10 CACHED

#11 [base 3/6] COPY package.json   ./
#11 CACHED

#12 [dev 3/4] RUN yarn
#12 CACHED

#13 [base 4/6] RUN yarn --production
#13 CACHED

#14 [base 5/6] RUN curl -sf https://gobinaries.com/tj/node-prune | sh
#14 CACHED

#15 [base 6/6] RUN node-prune
#15 CACHED

#16 [dev 2/4] COPY ./src/ ./src/
#16 CACHED

#17 [stage-2 3/4] COPY --from=dev /app/dist/ ./dist/
#17 CACHED

#18 [stage-2 2/4] COPY --from=base /app/package.json ./
#18 CACHED

#19 [stage-2 4/4] COPY --from=base /app/node_modules/ ./node_modules/
#19 CACHED

#20 exporting to image
#20 exporting layers done
#20 writing image sha256:476a556939f3a6243088edf27cd219e9b0980028f6f60aac7f09f1249691824d done
#20 naming to docker.io/library/672fbc21e58483169ee55b86 done
#20 DONE 0.0s

What's Next?
View summary of image vulnerabilities and recommendations â†’ docker scout quickview

docker image built successfully
identifying port
Detected service running port 3000
running your service
deployed successfully
 `;
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
