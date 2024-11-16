export interface DeploymentType {
  _id: string;
  title: string;
  latest_status: string;
  last_deployed_at: Date;
  repository_provider: string;
  domain_url: string;
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
  root_directory: string;
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

export interface PaginationType {
  total: number;
  current_page: number;
  last_page: number;
  per_page: number;
}

export interface UpdateDeploymentReqBody {
  title?: string;
  branch_name?: string;
  docker_file_path?: string;
  root_dir?: string | null;
}

export interface TResponse<T> {
  isSuccessful: boolean;
  data: T,
  error: string | null,
  code: number
}

export interface TPaginationResponse<T> {
  isSuccessful: boolean;
  data: T,
  error: string | null,
  code: number,
  pagination: PaginationType
}

export interface DeploymentEventType {
  id: string;
  deployment_id: string;
  title: string;
  type: string;
  triggered_by: string;
  triggered_value: string;
  status: string;
  reason: string;
  event_log_file_url: string;
  created_at: Date;
  updated_at: Date;
}
