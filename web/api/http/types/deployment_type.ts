export interface DeploymentType {
    _id: string;
    title: string;
    latest_status: string;
    last_deployed_at: Date
    repository_provider: string;
    branch_name: string;
    created_at: Date;
    updated_at: Date
}