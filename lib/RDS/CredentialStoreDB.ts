import { AmazonLinuxGeneration, AmazonLinuxImage, ISecurityGroup, IVpc, Instance, InstanceClass, InstanceSize, InstanceType, SubnetType } from "aws-cdk-lib/aws-ec2";
import { Credentials, DatabaseInstance, DatabaseInstanceEngine } from "aws-cdk-lib/aws-rds";
import { Construct } from "constructs";

export interface CredentialStoreDBProps {
  vpc: IVpc,
  databaseSecurityGroup: ISecurityGroup,
  ec2SecurityGroup: ISecurityGroup,
}

export class CredentialStoreDB extends Construct {

  readonly database : DatabaseInstance;

  constructor(scope: Construct, id:string, props: CredentialStoreDBProps) {
    super(scope, id);
    
    this.database = new DatabaseInstance(this, 'CredentialStoreDB', {
      databaseName: `CredentialStoreDB`,
      instanceIdentifier: `credentialstoredb`,
      engine: DatabaseInstanceEngine.POSTGRES,
      vpc: props.vpc,
      credentials: Credentials.fromGeneratedSecret('postgres', {
        secretName: 'credential-store-db-credentials'
      }),
      allocatedStorage: 20, // GiB
      cloudwatchLogsRetention: 14,
      instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
      publiclyAccessible: false,
      vpcSubnets: {
        subnetType: SubnetType.PRIVATE_ISOLATED
      },
      securityGroups: [props.databaseSecurityGroup]
    });
    
    new Instance(this, 'CredentialStoreDBEC2nstance', {
      instanceType: InstanceType.of(InstanceClass.T2, InstanceSize.MICRO),
      machineImage: new AmazonLinuxImage({ generation: AmazonLinuxGeneration.AMAZON_LINUX_2023 }),
      vpc: props.vpc,
      securityGroup: props.ec2SecurityGroup,
      associatePublicIpAddress: true,
      vpcSubnets: { subnetType: SubnetType.PUBLIC }
    })
  }
}