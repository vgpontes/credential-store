import { LambdaIntegration, RestApi } from "aws-cdk-lib/aws-apigateway";
import { Construct } from "constructs";
import { LambdaFunction } from "../Lambda/LambdaFunction";
import { Architecture, Runtime } from "aws-cdk-lib/aws-lambda";
import { RetentionDays } from "aws-cdk-lib/aws-logs";
import { ISecret } from "aws-cdk-lib/aws-secretsmanager";

export interface CredentialStoreApiGatewayProps {
    dbInfo : ISecret
}

export class CredentialStoreApiGateway extends Construct {
    constructor(scope: Construct, id: string, props: CredentialStoreApiGatewayProps) {
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
            codePath: './build/auth_service',
            description: 'API for creating an account, logging in to an account, and resetting password.',
            handler: 'handler',
            runtime: Runtime.PROVIDED_AL2023,
            architecture: Architecture.ARM_64,
            logGroupRetention: RetentionDays.TWO_WEEKS,
            environmentVariables: {
                "DB_SECRET_ID": props.dbInfo.secretName
            }
        });

        props.dbInfo.grantRead(authServiceLambda.lambdaFunction);

        const usersLambda = new LambdaFunction(this, 'UsersLambdaFn', {
            functionName: 'credential-store-users-api',
            codePath: './build/users',
            description: 'API for getting or adding users',
            handler: 'handler',
            runtime: Runtime.PROVIDED_AL2023,
            architecture: Architecture.ARM_64,
            logGroupRetention: RetentionDays.TWO_WEEKS
        });
      
        const authResource = api.root.addResource("authorize", { defaultIntegration: new LambdaIntegration(authServiceLambda.lambdaFunction) });
        authResource.addMethod("POST")

        const usersResource = api.root.addResource("users", { defaultIntegration: new LambdaIntegration(usersLambda.lambdaFunction) });
        usersResource.addMethod("GET");

        usersResource.addProxy();
    }
}