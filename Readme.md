# Deploybl

> Deploybl is a Platform as a Service (PaaS) that simplifies the deployment and management of applications.

## Table Of Contents

- [The goals of this project](#goal)
- [Features](#features)
- [Technologies](#technologies---libraries)
- [Documentation Apis](#documentation-apis)
- [Frontend](#frontend)

<a id="goal"></a>
## The Goals Of This Project
---

* Deploy any type of application (Node.js, Java, Go, etc.). No need to buy and configure a server.
* Generate a unique URL for each deployment. No need to buy and configure a domain.
* Using RabbitMQ as a Job Queue
* Building a simple dashboard with Next.js and shadcn/ui

<a id="features"></a>
### Features
---

* User Login/Registration
* Deploy any application from git repository
* Generate Unique URL for Application
* Manage Deployments
* View Deployment Logs

<a id="technologies"></a>
### Technologies
---

* [Go](https://golang.org/) - As Backend Language
* [Gin](https://github.com/gin-gonic/gin) - As Web Framework
* [Traefik](https://traefik.io/) - As Reverse Proxy
* [MongoDB](https://www.mongodb.com/) - As Database
* [Firebase](https://firebase.google.com/) - For Authentication
* [RabbitMQ](https://www.rabbitmq.com/) - As Job Queue
* [AWS S3](https://aws.amazon.com/s3/) - For log storage
* [Next.js](https://nextjs.org/) - For developing dashboard
* [shadcn/ui](https://ui.shadcn.com/) - UI compoenent
* [Tailwind CSS](https://tailwindcss.com/) - For styling
* [Docker](https://www.docker.com/) - For deployment
* [Github Actions](https://github.com/features/actions) - For CI

<a id="frontend"></a>
### Frontend
---
 <p> Signup Page </p>
<img src ="images/signup.jpeg">
 <p> Login Page </p>
<img src ="images/login.jpeg">
 <p> Dashboard Page </p>
<img src ="images/dashboard.jpeg">

 <p> Init New Deployment </p>
<img src ="images/init_new_deployment.jpeg">

 <p> Create New Deployment </p>
<img src ="images/create_new_deployment.jpeg">

 <p>Deployment Events</p>
<img src ="images/deployment_events.jpeg">

 <p> Deployment Event Logs </p>
<img src ="images/deployment_logs.jpeg">
 <p> Deployment Settings </p>
<img src ="images/deployment_settings.jpeg">
 <p> Deployment Env Update </p>
<img src ="images/env.jpeg">
