import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LoginService } from './LoginService/LoginService';
import { Vpc } from 'aws-cdk-lib/aws-ec2';

export class CredentialStoreStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id);

    const credentialStoreVpc = new Vpc(this, 'CredentialStoreVPC', {
      vpcName: 'credential-store-vpc',
    });

    new LoginService(this, 'CredentialStoreLoginService', {
        appName: 'credential-store',
        vpc: credentialStoreVpc
    });
  }
}
