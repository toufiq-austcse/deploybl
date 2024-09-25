export interface DeploymentType {
  _id: string;
  title: string;
  latest_status: string;
  last_deployed_at: Date;
  repository_provider: string;
  branch_name: string;
  created_at: Date;
  updated_at: Date;
}

export interface DeploymentDetailsType {
  _id: string;
  title: string;
  repository_name: string;
  domain_url: string;
  sub_domain_name: string;
  latest_status: string;
  last_deployed_at: Date;
  repository_provider: string;
  repository_url: string;
  branch_name: string;
  docker_file_path: string;
  env: object;
  created_at: Date;
  updated_at: Date;
}

export interface RepoDetailsType {
  svn_url: string;
  default_branch: string;
  name: string;
}

export interface CreateDeploymentApiRes {
  _id: string;
  title: string;
  latest_status: string;
  last_deployed_at: Date;
  repository_provider: string;
  branch_name: string;
  created_at: Date;
  updated_at: Date;
}

export interface DeploymentLatestStatusType {
  _id: string;
  latest_status: string;
  last_deployed_at: Date;
  domain_url: string;
}