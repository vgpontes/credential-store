import { LogGroup, RetentionDays } from "aws-cdk-lib/aws-logs";
import { Construct } from "constructs";
import { Architecture, Code, Function as LambdaFn, Runtime } from "aws-cdk-lib/aws-lambda";

export interface LambdaFunctionProps {
    functionName : string,
    description : string,
    codePath : string,
    handler : string,
    runtime : Runtime
    architecture : Architecture
    logGroupRetention : RetentionDays
    environmentVariables? : {[key:string]:string}
}

export class LambdaFunction extends Construct {

    readonly lambdaFunction : LambdaFn;

    constructor(scope: Construct, id: string, props: LambdaFunctionProps) {
        super(scope, id);

        const lambdaLogGroup = new LogGroup(this, `${props.functionName}-LambdaLogGroup`, {
            logGroupName: `/aws/lambda/${props.functionName}`,
            retention: props.logGroupRetention
        });
      
        this.lambdaFunction = new LambdaFn(this, `${props.functionName}-LambdaFn`, {
            functionName: props.functionName,
            code: Code.fromAsset(props.codePath),
            description: props.description,
            handler: props.handler,
            runtime: props.runtime,
            architecture: props.architecture,
            logGroup: lambdaLogGroup,
            environment: props.environmentVariables
        });
    }
}