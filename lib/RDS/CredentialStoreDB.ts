import { IVpc, InstanceClass, InstanceSize, InstanceType, SubnetType, Vpc } from "aws-cdk-lib/aws-ec2";
import { Credentials, DatabaseInstance, DatabaseInstanceEngine, IDatabaseInstance } from "aws-cdk-lib/aws-rds";
import { ISecret } from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from "constructs";

export class CredentialStoreDB extends Construct {

  readonly credentialStoreVpc : IVpc
  readonly database : DatabaseInstance;

  constructor(scope: Construct, id:string) {
    super(scope, id);
    
    this.credentialStoreVpc = new Vpc(this, 'CredentialStoreVPC', {
      vpcName: 'credential-store-vpc',
      enableDnsHostnames: true,
      enableDnsSupport: true,
      createInternetGateway: true,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'ingress',
          subnetType: SubnetType.PUBLIC
        },
        {
          cidrMask: 24,
          name: 'application',
          subnetType: SubnetType.PRIVATE_ISOLATED
        }
      ]
    });

    this.database = new DatabaseInstance(this, 'CredentialStoreUserDB', {
      databaseName: `CredentialStoreDB`,
      instanceIdentifier: `credentialstoredb`,
      engine: DatabaseInstanceEngine.POSTGRES,
      vpc: this.credentialStoreVpc,
      credentials: Credentials.fromGeneratedSecret('postgres', {
        secretName: 'credential-store-db-credentials'
      }),
      allocatedStorage: 20, // GiB
      cloudwatchLogsRetention: 14,
      instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
      publiclyAccessible: false
    });
  }
}