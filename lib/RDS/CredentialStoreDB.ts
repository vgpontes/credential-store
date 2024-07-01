import { InstanceClass, InstanceSize, InstanceType, Peer, Port, SecurityGroup, SubnetType, Vpc } from "aws-cdk-lib/aws-ec2";
import { Key } from "aws-cdk-lib/aws-kms";
import { Credentials, DatabaseInstance, DatabaseInstanceEngine } from "aws-cdk-lib/aws-rds";
import { ISecret, Secret } from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from "constructs";

export class CredentialStoreDB extends Construct {

  readonly dbInfo : ISecret;

  constructor(scope: Construct, id:string) {
    super(scope, id);
    
    const credentialStoreVpc = new Vpc(this, 'CredentialStoreVPC', {
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

    const database = new DatabaseInstance(this, 'CredentialStoreUserDB', {
      databaseName: `CredentialStoreDB`,
      instanceIdentifier: `credentialstoredb`,
      engine: DatabaseInstanceEngine.POSTGRES,
      vpc: credentialStoreVpc,
      credentials: Credentials.fromGeneratedSecret('postgres', {
        secretName: 'credential-store-db-credentials',
        
      }),
      allocatedStorage: 20, // GiB
      cloudwatchLogsRetention: 14,
      instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
      publiclyAccessible: true,
      vpcSubnets: {
        subnetType: SubnetType.PUBLIC
      }
    });

    database.connections.allowDefaultPortFromAnyIpv4();

    this.dbInfo = database.secret!;
  }
}