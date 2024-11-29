#!/bin/bash

# Cron Setup Script for Curl Commands

(crontab -l 2>/dev/null | grep -v "check-stopped-cron" | grep -v "check-deploying-cron") | crontab -
# Setup cron jobs
(crontab -l 2>/dev/null; echo "* * * * * curl --request POST --url 'http://localhost:3000/api/v1/deployments/check-stopped-cron'") | crontab -
(crontab -l 2>/dev/null; echo "* * * * * curl --request POST --url 'http://localhost:3000/api/v1/deployments/check-deploying-cron'") | crontab -

echo "Cron jobs for curl commands have been setup:"
