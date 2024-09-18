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
  sub_domain_name: string;
  latest_status: string;
  last_deployed_at: Date;
  repository_provider: string;
  repository_url: string;
  branch_name: string;
  docker_file_path: string;
  docker_image_tag: string;
  container_id: string;
  env: object;
  created_at: Date;
  updated_at: Date;
}