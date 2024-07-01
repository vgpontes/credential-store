import { App, Stack } from 'aws-cdk-lib';
import { InstanceClass, InstanceSize, InstanceType, Vpc } from 'aws-cdk-lib/aws-ec2';
import { DatabaseInstance, DatabaseInstanceEngine } from 'aws-cdk-lib/aws-rds';
import { CredentialStoreApiGateway } from './ApiGateway/ApiGateway';
import { CredentialStoreDB } from './RDS/Database';

export class CredentialStoreStack extends Stack {
  constructor(scope: App, id: string) {
    super(scope, id);

    const rdsDatabase = new CredentialStoreDB(this, 'CredentialStoreDB');

    new CredentialStoreApiGateway(this, 'CredentialStoreAPIGateway', {
      dbInfo: rdsDatabase.dbInfo
    });
  }
}
