import * as cdk from '@aws-cdk/core'
import * as lambda from '@aws-cdk/aws-lambda'
import * as iam from '@aws-cdk/aws-iam'
import * as api from '@aws-cdk/aws-apigateway'
import * as logs from '@aws-cdk/aws-logs'
import * as ddb from '@aws-cdk/aws-dynamodb'
import * as pure from 'aws-cdk-pure'
import * as hoc from 'aws-cdk-pure-hoc'
import * as path from 'path'

//
const app = new cdk.App()
const stack = new cdk.Stack(app, 'bookstore', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  }
})

//
const Storage = (): ddb.TableProps => ({
  partitionKey: {type: ddb.AttributeType.STRING, name: 'prefix'},
  sortKey: {type: ddb.AttributeType.STRING, name: 'suffix'},
  readCapacity: 1,
  writeCapacity: 1,
  removalPolicy: cdk.RemovalPolicy.DESTROY,
  tableName: `${cdk.Aws.STACK_NAME}-db`,
})
const storage = pure.join(stack, pure.iaac(ddb.Table)(Storage))


//
const Role = (): iam.RoleProps => ({
  assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
  managedPolicies: [
    iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole'),
  ],
  inlinePolicies: {
    'default': new iam.PolicyDocument({
      statements: [
        new iam.PolicyStatement({
          resources: [storage.tableArn],
          actions: [
            'dynamodb:*',
          ],
        })
      ]
    })
  }
})
const role = pure.join(stack, pure.iaac(iam.Role)(Role))

//
const Lambda = (): lambda.FunctionProps => ({
  code: hoc.common.AssetCodeGo(path.join(__dirname, '..')),
  handler: 'main',
  runtime: lambda.Runtime.GO_1_X,
  logRetention: logs.RetentionDays.FIVE_DAYS,
  functionName: `${cdk.Aws.STACK_NAME}-crud`,
  role,
  environment: {
    CONFIG_DDB: `ddb:///${storage.tableName}`
  }
})
const func = pure.wrap(api.LambdaIntegration)(pure.iaac(lambda.Function)(Lambda))

//
const Gateway = (): api.RestApiProps => ({
  deploy: true,
  deployOptions: { stageName: 'api' },
  endpointTypes: [api.EndpointType.REGIONAL],
  failOnWarnings: true,
  defaultCorsPreflightOptions: {
    allowOrigins: api.Cors.ALL_ORIGINS,
    maxAge: cdk.Duration.minutes(10),
  }
})
const rest = pure.iaac(api.RestApi)(Gateway)

//
pure.join(stack,
  pure.use({ rest, func })
  .effect(
    eff => {
      const seq = eff.rest.root.addResource('books')
      seq.addMethod('ANY', eff.func)

      const els = seq.addResource('{any+}')
      els.addMethod('ANY', eff.func)
    }
  )
)

app.synth()
