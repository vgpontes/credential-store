import { App, Stack } from 'aws-cdk-lib';
import { CredentialStoreApiGateway } from './ApiGateway/CredentialStoreApiGateway';
import { CredentialStoreDB } from './RDS/CredentialStoreDB';

export class CredentialStoreStack extends Stack {
  constructor(scope: App, id: string) {
    super(scope, id);

    const rdsDatabase = new CredentialStoreDB(this, 'CredentialStoreDB');

    new CredentialStoreApiGateway(this, 'CredentialStoreAPIGateway', {
      database: rdsDatabase.database,
      vpc: rdsDatabase.credentialStoreVpc
    });
  }
}
