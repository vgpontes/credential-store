import { InstanceClass, InstanceSize, InstanceType, Vpc } from "aws-cdk-lib/aws-ec2";
import { DatabaseInstance, DatabaseInstanceEngine } from "aws-cdk-lib/aws-rds";
import { Construct } from "constructs";

export class CredentialStoreDB extends Construct {
    constructor(scope: Construct, id:string) {
        super(scope, id);

        const credentialStoreVpc = new Vpc(this, 'CredentialStoreVPC', {
            vpcName: 'credential-store-vpc',
          });
      
          new DatabaseInstance(this, 'CredentialStoreUserDB', {
            databaseName: `CredentialStoreDB`,
            instanceIdentifier: `credentialstoredb`,
            engine: DatabaseInstanceEngine.POSTGRES,
            vpc: credentialStoreVpc,
            allocatedStorage: 20, // GiB
            cloudwatchLogsRetention: 14,
            instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
            publiclyAccessible: false
          });
    }
}