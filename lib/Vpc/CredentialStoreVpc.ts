import { Construct } from "constructs";
import { ISecurityGroup, IVpc, Port, SecurityGroup, SubnetType, Vpc } from "aws-cdk-lib/aws-ec2";

export class CredentialStoreVpc extends Construct {

    readonly credentialStoreVpc : IVpc;
    readonly rdsSecurityGroup : ISecurityGroup;
    readonly ec2SecurityGroup : ISecurityGroup;
    readonly lambdaSecurityGroup : ISecurityGroup;

    constructor(scope: Construct, id: string) {
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
      
        this.rdsSecurityGroup = new SecurityGroup(this, 'RDSSecurityGroup', {
            securityGroupName: 'rds-ec2-0',
            vpc: this.credentialStoreVpc,
            allowAllOutbound: false,
        })
      
        this.ec2SecurityGroup = new SecurityGroup(this, 'EC2SecurityGroup', {
            securityGroupName: 'ec2-rds-0',
            vpc: this.credentialStoreVpc,
            allowAllOutbound: false,
        })
      
        const ec2ConnectEndpointSecurityGroup = new SecurityGroup(this, 'EC2InstanceConnectEndpointSecurityGroup', {
            securityGroupName: 'ec2-instance-connect-endpoint-sg',
            vpc: this.credentialStoreVpc,
            allowAllOutbound: false,
        })

        this.lambdaSecurityGroup = new SecurityGroup(this, 'LambdaDatabaseSecurityGroup', {
            securityGroupName: 'lambda-rds-0',
            vpc: this.credentialStoreVpc,
            allowAllOutbound: false,
        })
      
        this.rdsSecurityGroup.addIngressRule(this.ec2SecurityGroup, Port.POSTGRES, 'Ingress for EC2 access');
        this.rdsSecurityGroup.addIngressRule(this.lambdaSecurityGroup, Port.POSTGRES, 'Ingress for Lambda access')
        this.ec2SecurityGroup.addEgressRule(this.rdsSecurityGroup, Port.POSTGRES, 'Egress to access RDS database');
        this.ec2SecurityGroup.addIngressRule(ec2ConnectEndpointSecurityGroup, Port.tcp(22), 'Allows inbound SSH traffic from the resources associated with the endpoint security group');
        ec2ConnectEndpointSecurityGroup.addEgressRule(this.ec2SecurityGroup, Port.tcp(22), 'Allows outbound SSH traffic to all instances associated with the instance security group');
        this.lambdaSecurityGroup.addEgressRule(this.rdsSecurityGroup, Port.POSTGRES, 'Egress to access RDS database');
    }
}