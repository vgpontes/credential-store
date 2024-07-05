import { App, Stack } from 'aws-cdk-lib';
import { CredentialStoreApiGateway } from './ApiGateway/CredentialStoreApiGateway';
import { CredentialStoreDB } from './RDS/CredentialStoreDB';
import { CredentialStoreVpc } from './Vpc/CredentialStoreVpc';

export class CredentialStoreStack extends Stack {
  constructor(scope: App, id: string) {
    super(scope, id);

    const vpc = new CredentialStoreVpc(this, 'CredentialStoreVPC');

    const rdsDatabase = new CredentialStoreDB(this, 'CredentialStoreDB', {
      vpc: vpc.credentialStoreVpc,
      ec2SecurityGroup: vpc.ec2SecurityGroup,
      databaseSecurityGroup: vpc.rdsSecurityGroup
    });

    new CredentialStoreApiGateway(this, 'CredentialStoreAPIGateway', {
      database: rdsDatabase.database,
      vpc: vpc.credentialStoreVpc,
      lambdaSecurityGroup: vpc.lambdaSecurityGroup
    });
  }
}
