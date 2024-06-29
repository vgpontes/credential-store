import { Construct } from 'constructs';
import { LambdaRestApi, ResponseType } from 'aws-cdk-lib/aws-apigateway';
import { Function, Code, Runtime, Architecture } from 'aws-cdk-lib/aws-lambda';
import { LogGroup, RetentionDays } from 'aws-cdk-lib/aws-logs';
import { DatabaseInstance, DatabaseInstanceEngine } from 'aws-cdk-lib/aws-rds';
import { IVpc, InstanceClass, InstanceSize, InstanceType } from 'aws-cdk-lib/aws-ec2';

export interface LoginServiceProps {
  appName: string,
  vpc: IVpc
}

export class LoginService extends Construct {
  constructor(scope: Construct, id: string, props: LoginServiceProps) {
    super(scope, id);

    const lambdaLogGroup = new LogGroup(this, 'LoginServiceLambdaLogGroup', {
      logGroupName: `/aws/lambda/${props.appName}-login-service`,
      retention: RetentionDays.TWO_WEEKS
    });

    const lambdaFunction = new Function(this, 'LoginServiceLambdaFn', {
      functionName: `${props.appName}-login-service`,
      code: Code.fromAsset('./build/loginService'),
      description: 'API for creating an account, logging in to an account, and resetting password.',
      handler: 'handler',
      runtime: Runtime.PROVIDED_AL2023,
      architecture: Architecture.ARM_64,
      logGroup: lambdaLogGroup
    });

    const loginApi = new LambdaRestApi(this, 'LoginServiceRestApi', {
      handler: lambdaFunction,
      description: `REST API for ${props.appName} Login Service.`,
    });

    loginApi.addGatewayResponse('LoginServiceUnauthenticated', {
      type: ResponseType.MISSING_AUTHENTICATION_TOKEN,
      templates: {
        'application/json': '{ "message": $context.error.messageString, "statusCode": "488", "type": "$context.error.responseType" }'
      }
    });

    loginApi.addGatewayResponse('LoginServiceUnauthorized', {
      type: ResponseType.ACCESS_DENIED,
      templates: {
        'application/json': '{ "message": $context.error.messageString, "statusCode": "488", "type": "$context.error.responseType" }'
      }
    });

    new DatabaseInstance(this, 'UserDB', {
      databaseName: `${props.appName}-users`,
      instanceIdentifier: `${props.appName}-users`,
      engine: DatabaseInstanceEngine.POSTGRES,
      vpc: props.vpc,
      allocatedStorage: 20, // GiB
      cloudwatchLogsRetention: 14,
      instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
      publiclyAccessible: false
    });
  }
}
