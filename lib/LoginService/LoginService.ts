import { Size } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LambdaRestApi } from 'aws-cdk-lib/aws-apigateway';
import { Function, Code, Runtime } from 'aws-cdk-lib/aws-lambda';
import { LogGroup, RetentionDays } from 'aws-cdk-lib/aws-logs';

export interface LoginServiceProps {
  appName: string,
}

export class LoginService extends Construct {
  constructor(scope: Construct, id: string, props: LoginServiceProps) {
    super(scope, id);

    const lambdaLogGroup = new LogGroup(this, 'LoginServiceLambdaLogGroup', {
      logGroupName: `/aws/lambda/${props.appName}`,
      retention: RetentionDays.TWO_WEEKS
    });

    const lambdaFunction = new Function(this, 'LoginServiceLambdaFn', {
      functionName: `${props.appName}-login-service`,
      code: Code.fromAsset('./src/loginService'),
      description: 'API for creating an account, logging in to an account, and resetting password.',
      handler: 'handler',
      runtime: Runtime.PROVIDED_AL2023,
      logGroup: lambdaLogGroup
    });

    new LambdaRestApi(this, 'LoginServiceRestApi', {
      handler: lambdaFunction,
      description: `REST API for ${props.appName} Login Service.`,
    });
  }
}
