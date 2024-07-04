import { AmazonLinux2023ImageSsmParameter, AmazonLinuxGeneration, AmazonLinuxImage, CfnInstanceConnectEndpoint, IVpc, Instance, InstanceClass, InstanceSize, InstanceType, Port, SecurityGroup, SubnetType, Vpc } from "aws-cdk-lib/aws-ec2";
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

    const rdsSecurityGroup = new SecurityGroup(this, 'RDSSecurityGroup', {
      securityGroupName: 'rds-ec2-0',
      vpc: this.credentialStoreVpc,
      allowAllOutbound: false,
    })

    const ec2SecurityGroup = new SecurityGroup(this, 'EC2SecurityGroup', {
      securityGroupName: 'ec2-rds-0',
      vpc: this.credentialStoreVpc,
      allowAllOutbound: false,
    })

    const ec2ConnectEndpointSecurityGroup = new SecurityGroup(this, 'EC2InstanceConnectEndpointSecurityGroup', {
      securityGroupName: 'ec2-instance-connect-endpoint-sg',
      vpc: this.credentialStoreVpc,
      allowAllOutbound: false,
    })

    rdsSecurityGroup.addIngressRule(ec2SecurityGroup, Port.POSTGRES, 'Ingress for RDS instance');
    ec2SecurityGroup.addEgressRule(rdsSecurityGroup, Port.POSTGRES, 'Egress for EC2 instance');
    ec2SecurityGroup.addIngressRule(ec2ConnectEndpointSecurityGroup, Port.tcp(22), 'Allows inbound SSH traffic from the resources associated with the endpoint security group')
    ec2ConnectEndpointSecurityGroup.addEgressRule(ec2SecurityGroup, Port.tcp(22), 'Allows outbound SSH traffic to all instances associated with the instance security group');
    
    this.database = new DatabaseInstance(this, 'CredentialStoreDB', {
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
      publiclyAccessible: false,
      vpcSubnets: {
        subnetType: SubnetType.PRIVATE_ISOLATED
      },
      securityGroups: [rdsSecurityGroup]
    });

    new Instance(this, 'CredentialStoreDBEC2nstance', {
      instanceType: InstanceType.of(InstanceClass.T2, InstanceSize.MICRO),
      machineImage: new AmazonLinuxImage({ generation: AmazonLinuxGeneration.AMAZON_LINUX_2023 }),
      vpc: this.credentialStoreVpc,
      securityGroup: ec2SecurityGroup
    })

    new CfnInstanceConnectEndpoint(this, 'EC2InstanceConnectEndpoint', {
      subnetId: this.credentialStoreVpc.isolatedSubnets[0].subnetId,
      securityGroupIds: [ec2ConnectEndpointSecurityGroup.uniqueId]
    })
  }
}