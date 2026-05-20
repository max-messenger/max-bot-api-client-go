AWS()                                                                    AWS()

NNAAMMEE
       aws -

DDEESSCCRRIIPPTTIIOONN
       The AWS Command Line Interface is a unified tool to manage your AWS
       services.

SSYYNNOOPPSSIISS

          aws [options] <command> <subcommand> [parameters]

       Use _a_w_s _c_o_m_m_a_n_d _h_e_l_p for information on a specific command. Use _a_w_s
       _h_e_l_p _t_o_p_i_c_s to view a list of available help topics. The synopsis for
       each command shows its parameters and their usage. Optional parameters
       are shown in square brackets.

GGLLOOBBAALL OOPPTTIIOONNSS
       ----ddeebbuugg (boolean)

       Turn on debug logging.

       ----eennddppooiinntt--uurrll (string)

       Override command's default URL with the given URL.

       ----nnoo--vveerriiffyy--ssssll (boolean)

       By default, the AWS CLI uses SSL when communicating with AWS services.
       For each SSL connection, the AWS CLI will verify SSL certificates. This
       option overrides the default behavior of verifying SSL certificates.

       ----nnoo--ppaaggiinnaattee (boolean)

       Disable automatic pagination. If automatic pagination is disabled, the
       AWS CLI will only make one call, for the first page of results.

       ----oouuttppuutt (string)

       The formatting style for command output.

       +o json

       +o text

       +o table

       +o yaml

       +o yaml-stream

       +o off

       ----qquueerryy (string)

       A JMESPath query to use in filtering the response data.

       ----pprrooffiillee (string)

       Use a specific profile from your credential file.

       ----rreeggiioonn (string)

       The region to use. Overrides config/env settings.

       ----vveerrssiioonn (string)

       Display the version of this tool.

       ----ccoolloorr (string)

       Turn on/off color output.

       +o on

       +o off

       +o auto

       ----nnoo--ssiiggnn--rreeqquueesstt (boolean)

       Do not sign requests. Credentials will not be loaded if this argument
       is provided.

       ----ccaa--bbuunnddllee (string)

       The CA certificate bundle to use when verifying SSL certificates.
       Overrides config/env settings.

       ----ccllii--rreeaadd--ttiimmeeoouutt (int)

       The maximum socket read time in seconds. If the value is set to 0, the
       socket read will be blocking and not timeout. The default value is 60
       seconds.

       ----ccllii--ccoonnnneecctt--ttiimmeeoouutt (int)

       The maximum socket connect time in seconds. If the value is set to 0,
       the socket connect will be blocking and not timeout. The default value
       is 60 seconds.

       ----ccllii--bbiinnaarryy--ffoorrmmaatt (string)

       The formatting style to be used for binary blobs. The default format is
       base64. The base64 format expects binary blobs to be provided as a
       base64 encoded string. The raw-in-base64-out format preserves
       compatibility with AWS CLI V1 behavior and binary values must be passed
       literally. When providing contents from a file that map to a binary
       blob ffiilleebb:://// will always be treated as binary and use the file
       contents directly regardless of the ccllii--bbiinnaarryy--ffoorrmmaatt setting. When
       using ffiillee:://// the file contents will need to properly formatted for the
       configured ccllii--bbiinnaarryy--ffoorrmmaatt.

       +o base64

       +o raw-in-base64-out

       ----nnoo--ccllii--ppaaggeerr (boolean)

       Disable cli pager for output.

       ----ccllii--aauuttoo--pprroommpptt (boolean)

       Automatically prompt for CLI input parameters.

       ----nnoo--ccllii--aauuttoo--pprroommpptt (boolean)

       Disable automatically prompt for CLI input parameters.

       ----ccllii--eerrrroorr--ffoorrmmaatt (string)

       The formatting style for error output. By default, errors are displayed
       in enhanced format.

       +o legacy

       +o json

       +o yaml

       +o text

       +o table

       +o enhanced

AAVVAAIILLAABBLLEE SSEERRVVIICCEESS

       +o accessanalyzer

       +o account

       +o acm

       +o acm-pca

       +o aiops

       +o amp

       +o amplify

       +o amplifybackend

       +o amplifyuibuilder

       +o apigateway

       +o apigatewaymanagementapi

       +o apigatewayv2

       +o appconfig

       +o appconfigdata

       +o appfabric

       +o appflow

       +o appintegrations

       +o application-autoscaling

       +o application-insights

       +o application-signals

       +o applicationcostprofiler

       +o appmesh

       +o apprunner

       +o appstream

       +o appsync

       +o arc-region-switch

       +o arc-zonal-shift

       +o artifact

       +o athena

       +o auditmanager

       +o autoscaling

       +o autoscaling-plans

       +o b2bi

       +o backup

       +o backup-gateway

       +o backupsearch

       +o batch

       +o bcm-dashboards

       +o bcm-data-exports

       +o bcm-pricing-calculator

       +o bcm-recommended-actions

       +o bedrock

       +o bedrock-agent

       +o bedrock-agent-runtime

       +o bedrock-agentcore

       +o bedrock-agentcore-control

       +o bedrock-data-automation

       +o bedrock-data-automation-runtime

       +o bedrock-runtime

       +o billing

       +o billingconductor

       +o braket

       +o budgets

       +o ce

       +o chatbot

       +o chime

       +o chime-sdk-identity

       +o chime-sdk-media-pipelines

       +o chime-sdk-meetings

       +o chime-sdk-messaging

       +o chime-sdk-voice

       +o cleanrooms

       +o cleanroomsml

       +o cli-dev

       +o cloud9

       +o cloudcontrol

       +o clouddirectory

       +o cloudformation

       +o cloudfront

       +o cloudfront-keyvaluestore

       +o cloudhsm

       +o cloudhsmv2

       +o cloudsearch

       +o cloudsearchdomain

       +o cloudtrail

       +o cloudtrail-data

       +o cloudwatch

       +o codeartifact

       +o codebuild

       +o codecatalyst

       +o codecommit

       +o codeconnections

       +o codeguru-reviewer

       +o codeguru-security

       +o codeguruprofiler

       +o codepipeline

       +o codestar-connections

       +o codestar-notifications

       +o cognito-identity

       +o cognito-idp

       +o cognito-sync

       +o comprehend

       +o comprehendmedical

       +o compute-optimizer

       +o compute-optimizer-automation

       +o configservice

       +o configure

       +o connect

       +o connect-contact-lens

       +o connectcampaigns

       +o connectcampaignsv2

       +o connectcases

       +o connecthealth

       +o connectparticipant

       +o controlcatalog

       +o controltower

       +o cost-optimization-hub

       +o cur

       +o customer-profiles

       +o databrew

       +o dataexchange

       +o datapipeline

       +o datasync

       +o datazone

       +o dax

       +o ddb

       +o deadline

       +o deploy

       +o detective

       +o devicefarm

       +o devops-agent

       +o devops-guru

       +o directconnect

       +o discovery

       +o dlm

       +o dms

       +o docdb

       +o docdb-elastic

       +o drs

       +o ds

       +o ds-data

       +o dsql

       +o dynamodb

       +o dynamodbstreams

       +o ebs

       +o ec2

       +o ec2-instance-connect

       +o ecr

       +o ecr-public

       +o ecs

       +o efs

       +o eks

       +o eks-auth

       +o elasticache

       +o elasticbeanstalk

       +o elb

       +o elbv2

       +o elementalinference

       +o emr

       +o emr-containers

       +o emr-serverless

       +o entityresolution

       +o es

       +o events

       +o evs

       +o finspace

       +o finspace-data

       +o firehose

       +o fis

       +o fms

       +o forecast

       +o forecastquery

       +o frauddetector

       +o freetier

       +o fsx

       +o gamelift

       +o gameliftstreams

       +o geo-maps

       +o geo-places

       +o geo-routes

       +o glacier

       +o globalaccelerator

       +o glue

       +o grafana

       +o greengrass

       +o greengrassv2

       +o groundstation

       +o guardduty

       +o health

       +o healthlake

       +o help

       +o history

       +o iam

       +o identitystore

       +o imagebuilder

       +o importexport

       +o inspector

       +o inspector-scan

       +o inspector2

       +o interconnect

       +o internetmonitor

       +o invoicing

       +o iot

       +o iot-data

       +o iot-jobs-data

       +o iot-managed-integrations

       +o iotdeviceadvisor

       +o iotevents

       +o iotevents-data

       +o iotfleetwise

       +o iotsecuretunneling

       +o iotsitewise

       +o iotthingsgraph

       +o iottwinmaker

       +o iotwireless

       +o ivs

       +o ivs-realtime

       +o ivschat

       +o kafka

       +o kafkaconnect

       +o kendra

       +o kendra-ranking

       +o keyspaces

       +o keyspacesstreams

       +o kinesis

       +o kinesis-video-archived-media

       +o kinesis-video-media

       +o kinesis-video-signaling

       +o kinesis-video-webrtc-storage

       +o kinesisanalytics

       +o kinesisanalyticsv2

       +o kinesisvideo

       +o kms

       +o lakeformation

       +o lambda

       +o launch-wizard

       +o lex-models

       +o lex-runtime

       +o lexv2-models

       +o lexv2-runtime

       +o license-manager

       +o license-manager-linux-subscriptions

       +o license-manager-user-subscriptions

       +o lightsail

       +o location

       +o login

       +o logout

       +o logs

       +o lookoutequipment

       +o m2

       +o machinelearning

       +o macie2

       +o mailmanager

       +o managedblockchain

       +o managedblockchain-query

       +o marketplace-agreement

       +o marketplace-catalog

       +o marketplace-deployment

       +o marketplace-discovery

       +o marketplace-entitlement

       +o marketplace-reporting

       +o marketplacecommerceanalytics

       +o mediaconnect

       +o mediaconvert

       +o medialive

       +o mediapackage

       +o mediapackage-vod

       +o mediapackagev2

       +o mediastore

       +o mediastore-data

       +o mediatailor

       +o medical-imaging

       +o memorydb

       +o meteringmarketplace

       +o mgh

       +o mgn

       +o migration-hub-refactor-spaces

       +o migrationhub-config

       +o migrationhuborchestrator

       +o migrationhubstrategy

       +o mpa

       +o mq

       +o mturk

       +o mwaa

       +o mwaa-serverless

       +o neptune

       +o neptune-graph

       +o neptunedata

       +o network-firewall

       +o networkflowmonitor

       +o networkmanager

       +o networkmonitor

       +o notifications

       +o notificationscontacts

       +o nova-act

       +o oam

       +o observabilityadmin

       +o odb

       +o omics

       +o opensearch

       +o opensearchserverless

       +o organizations

       +o osis

       +o outposts

       +o panorama

       +o partnercentral-account

       +o partnercentral-benefits

       +o partnercentral-channel

       +o partnercentral-selling

       +o payment-cryptography

       +o payment-cryptography-data

       +o pca-connector-ad

       +o pca-connector-scep

       +o pcs

       +o personalize

       +o personalize-events

       +o personalize-runtime

       +o pi

       +o pinpoint

       +o pinpoint-email

       +o pinpoint-sms-voice

       +o pinpoint-sms-voice-v2

       +o pipes

       +o polly

       +o pricing

       +o proton

       +o qapps

       +o qbusiness

       +o qconnect

       +o quicksight

       +o ram

       +o rbin

       +o rds

       +o rds-data

       +o redshift

       +o redshift-data

       +o redshift-serverless

       +o rekognition

       +o repostspace

       +o resiliencehub

       +o resource-explorer-2

       +o resource-groups

       +o resourcegroupstaggingapi

       +o rolesanywhere

       +o route53

       +o route53-recovery-cluster

       +o route53-recovery-control-config

       +o route53-recovery-readiness

       +o route53domains

       +o route53globalresolver

       +o route53profiles

       +o route53resolver

       +o rtbfabric

       +o rum

       +o s3

       +o s3api

       +o s3control

       +o s3files

       +o s3outposts

       +o s3tables

       +o s3vectors

       +o sagemaker

       +o sagemaker-a2i-runtime

       +o sagemaker-edge

       +o sagemaker-featurestore-runtime

       +o sagemaker-geospatial

       +o sagemaker-metrics

       +o sagemaker-runtime

       +o savingsplans

       +o scheduler

       +o schemas

       +o sdb

       +o secretsmanager

       +o security-ir

       +o securityagent

       +o securityhub

       +o securitylake

       +o serverlessrepo

       +o service-quotas

       +o servicecatalog

       +o servicecatalog-appregistry

       +o servicediscovery

       +o ses

       +o sesv2

       +o shield

       +o signer

       +o signer-data

       +o signin

       +o simpledbv2

       +o simspaceweaver

       +o snow-device-management

       +o snowball

       +o sns

       +o socialmessaging

       +o sqs

       +o ssm

       +o ssm-contacts

       +o ssm-guiconnect

       +o ssm-incidents

       +o ssm-quicksetup

       +o ssm-sap

       +o sso

       +o sso-admin

       +o sso-oidc

       +o stepfunctions

       +o storagegateway

       +o sts

       +o supplychain

       +o support

       +o support-app

       +o sustainability

       +o swf

       +o synthetics

       +o taxsettings

       +o textract

       +o timestream-influxdb

       +o timestream-query

       +o timestream-write

       +o tnb

       +o transcribe

       +o transfer

       +o translate

       +o trustedadvisor

       +o uxc

       +o verifiedpermissions

       +o voice-id

       +o vpc-lattice

       +o waf

       +o waf-regional

       +o wafv2

       +o wellarchitected

       +o wickr

       +o wisdom

       +o workdocs

       +o workmail

       +o workmailmessageflow

       +o workspaces

       +o workspaces-instances

       +o workspaces-thin-client

       +o workspaces-web

       +o xray

SSEEEE AALLSSOO

       +o aws help topics

                                                                         AWS()
