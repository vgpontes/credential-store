#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { CredentialStoreStack } from '../lib/CredentialStoreStack';

const app = new cdk.App();

new CredentialStoreStack(app, 'CredentialStore');