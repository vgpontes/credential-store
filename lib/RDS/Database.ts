import { InstanceClass, InstanceSize, InstanceType, Vpc } from "aws-cdk-lib/aws-ec2";
import { DatabaseInstance, DatabaseInstanceEngine } from "aws-cdk-lib/aws-rds";
import { ISecret } from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from "constructs";

export class CredentialStoreDB extends Construct {

  readonly dbInfo : ISecret;

  constructor(scope: Construct, id:string) {
    super(scope, id);

    const credentialStoreVpc = new Vpc(this, 'CredentialStoreVPC', {
      vpcName: 'credential-store-vpc',
    });
      
    const database = new DatabaseInstance(this, 'CredentialStoreUserDB', {
      databaseName: `CredentialStoreDB`,
      instanceIdentifier: `credentialstoredb`,
      engine: DatabaseInstanceEngine.POSTGRES,
      vpc: credentialStoreVpc,
      allocatedStorage: 20, // GiB
      cloudwatchLogsRetention: 14,
      instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
      publiclyAccessible: true
    });

    this.dbInfo = database.secret?
  }
}