import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LoginService } from './LoginService/LoginService';

export class CredentialStoreStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id);

    new LoginService(this, 'CredentialStoreLoginService', {
        appName: 'credential-store'
    });
  }
}
