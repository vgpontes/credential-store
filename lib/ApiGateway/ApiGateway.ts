import { LambdaIntegration, RestApi } from "aws-cdk-lib/aws-apigateway";
import { Construct } from "constructs";
import { LambdaFunction } from "../Lambda/AuthService";
import { Architecture, Runtime } from "aws-cdk-lib/aws-lambda";
import { RetentionDays } from "aws-cdk-lib/aws-logs";

export class CredentialStoreApiGateway extends Construct {
    constructor(scope: Construct, id: string) {
        super(scope, id);
      
        const api = new RestApi(this, 'CredentialStoreApi', {
            restApiName: 'Credential-Store API',
            description: "API Gateway for Credential-Store",
            deployOptions: {
                stageName: "api"
            }
        });

        const authServiceLambda = new LambdaFunction(this, 'AuthServiceLambdaFn', {
            functionName: 'credential-store-auth-service',
            codePath: './build/authService',
            description: 'API for creating an account, logging in to an account, and resetting password.',
            handler: 'handler',
            runtime: Runtime.PROVIDED_AL2023,
            architecture: Architecture.ARM_64,
            logGroupRetention: RetentionDays.TWO_WEEKS
        })
      
        const authResource = api.root.addResource("authorize");
        authResource.addMethod("POST", new LambdaIntegration(authServiceLambda.lambdaFunction))
    }
}