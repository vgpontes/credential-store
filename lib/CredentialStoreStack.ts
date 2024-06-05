import { Stack } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LoginService } from './LoginService/LoginService';

export class CredentialStoreStack extends Stack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new LoginService(this, 'CredentialStoreLoginService', {
        appName: 'credential-store'
    });
  }
}
