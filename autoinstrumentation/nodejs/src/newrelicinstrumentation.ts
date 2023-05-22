const newrelic = require('newrelic');
const awsSdk = require('@newrelic/aws-sdk');
const koa = require('@newrelic/koa');
const superagent = require('@newrelic/superagent');
const nativeMetrics = require('@newrelic/native-metrics');

export { newrelic, awsSdk, koa, superagent, nativeMetrics };