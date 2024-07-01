import { Construct } from "constructs";
import { IVpc, PublicSubnet, Vpc } from "aws-cdk-lib/aws-ec2";

export class CredentialStoreVpc extends Construct {

    readonly vpc : IVpc;

    constructor(scope: Construct, id: string) {
        super(scope, id);
      
        this.vpc = new Vpc(this, 'CredentialStoreVPC', {
            vpcName: 'credential-store-vpc',
            enableDnsHostnames: true,
            enableDnsSupport: true
        });

        new PublicSubnet(this, 'CredentialStorePublicSubnet', {
            vpcId: this.vpc.vpcId,
            cidrBlock: '',
            availabilityZone: ''
        })
    }
}