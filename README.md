
Setting up Performance Test Framework
=========================

This documents aims to walk you through the setup required to run the performance test rig for **Credhub** on your own.

----------
## Table of Contents
1. [IaaS Setup for the BOSH Deployment](#iaas-setup-for-the-bosh-deployment)
1. [Deploying CredHub via Bosh](#deploying-credhub-via-bosh)
1. [Running the Performance Test](#running-the-performance-test)

## IaaS Setup for the BOSH Deployment

The subject is deployed entirely on AWS. We are testing a Credhub cluster and each Credhub is deployed on a m4.large VM behind a load-balancer. The UAA is also deployed on a separate m4.large VM. The Credhub cluster uses a Postgres RDS as its datastore while the UAA uses an internally deployed postgres.

#### Setting up a load balancer 

Update your *cloudformation.json* with the the configuration required to add a load balancer which resembles the following:

```json
"PerformanceLoadBalancer": {
      "Type": "AWS::ElasticLoadBalancing::LoadBalancer",
      "DependsOn": "AttachGateway",
      "Properties": {
        "Listeners": [{
          "LoadBalancerPort": "8844",
          "InstancePort": "8844",
          "Protocol": "TCP",
          "InstanceProtocol": "TCP"
        }],
        "SecurityGroups": [{"Ref": "LoadBalancerSecurityGroup"}],
        "Subnets": [
          {"Ref": "YourPrivateSubnet"}
        ],
        "HealthCheck": {
          "Target": "TCP:8844",
          "HealthyThreshold": "3",
          "UnhealthyThreshold": "5",
          "Interval": "30",
          "Timeout": "5"
        }
      }
    }
```

Upload the updated *cloudformation.json* to AWS which should create your load-balancer.

On the AWS console,  under the EC2 Dashboard, you will have a new addition in the LoadBalancers section. 
This would be a good time to create a CNAME in the Route53 Dashboard for the **DNS Name** of your LoadBalancer.
Additionally note the **Name** of the LoadBalancer as this will be required while defining a vm_type for your Credhub instances.

Now in the *cloud-config.yml* file, add a section under *vm_types* as follows:

```yml
- name: performance
  cloud_properties:
    instance_type: m4.large
    elbs: ["LoadBalancerName"]
    auto_assign_public_ip: true
    ephemeral_disk:
      size: 25000
      type: gp2
```

> **Note:**
Particular attention must be placed to the *auto_assign_public_ip* flag. In our configuration we ensure UAA can accept incoming traffic from all IPs. Ensure your security groups are defined as such. In case you want to ensure only requests from certain IPs are forwarded to UAA, ensure the *auto_assign_public_ip* flag is false, and ensure an elastic IP is assigned to each Credhub instance and that communication from each of those IPs are allowed by the security group UAA has been defined in.

#### Setting up an RDS instance

Find below the addition required to *cloudformation.json* to setup an RDS instance with Postgres running on it. You will need to create a DBSubnetGroup for it and ensure that the DBSubnetGroup has subnets within it that are on different availability zones. 

```json
    "CredHubPerformanceRDSSubnet": {
      "Type" : "AWS::RDS::DBSubnetGroup",
      "Properties" : {
        "DBSubnetGroupDescription" : "Subnet for performance RDS instance",
        "SubnetIds" : [
          { "Ref": "CredHubBoshSubnet" },
          { "Ref": "CredHubBoshSubnetTwo" }
        ]
      }
    }
```
```json
    "CredHubPerformanceRDS": {
      "Type": "AWS::RDS::DBInstance",
      "DependsOn": "CredHubPerformanceRDSSubnet",
      "Properties": {
        "AllocatedStorage": "50",
        "AllowMajorVersionUpgrade": false,
        "AutoMinorVersionUpgrade": true,
        "AvailabilityZone": "us-east-1a",
        "DBInstanceClass": "db.m4.2xlarge",
        "DBInstanceIdentifier": "credhubPerf",
        "DBSubnetGroupName": {
          "Ref": "CredHubPerformanceRDSSubnet"
        },
        "DBName": "your-db-name",
        "Engine": "postgres",
        "EngineVersion": "9.4.11",
        "MasterUsername": "postgres",
        "MasterUserPassword": "your-db-password",
        "MultiAZ": false,
        "Port": "5432",
        "PubliclyAccessible": false,
        "StorageEncrypted": false,
        "StorageType": "gp2",
        "VPCSecurityGroups": [
          {
            "Ref": "CredHubSecurityGroup"
          }
        ]
      }
    }

```

Ensure you have updated your cloud-config using:
```bash
bosh -e boshenv update-cloud-config /path/to/cloud_config.yml
``` 

## Deploying CredHub via Bosh

1. Create a `vars` file with properties specific to your CredHub deployment
   ```yml
   credhub_ca_name: #REPLACE ME
   db_password: #REPLACE ME
   db_host: #REPLACE ME
   db_port: #REPLACE ME
   db_name: #REPLACE ME
   uaa_url: #REPLACE ME
   uaa_verification_key: #REPLACE ME
   uaa_ca: #REPLACE ME
   ```
1. Deploy credhub with `credhub_cannon` errand co-located using `sample-manifests/credhub-ha-perf.yml`

   ```bash
   bosh deploy sample-manifests/credhub_ha_perf.yml \
    -v instances=<INSTANCES> \
    -v min_concurrent=<MIN_CONCURRENT> \
    -v max_concurrent=<MAX_CONCURRENT> \
    -v step_size=<STEP_SIZE> \
    -v request_type=<REQUEST_TYPE> \
    -v num_requests=<NUM_REQUESTS> \
    -v credhub_host=<CREDHUB_HOST> \
    -l /path/to/vars/file
   ```

Now you wait.

> **Troubleshooting common deployment failures:**
> - Check your firewall rules. Are  they setup the way you expect them to be. Can the Credhub instances talk to UAA.
> - Check your certificates. Do the CA's have the rights IPs defined in them? Are you passing around the correct certificates?

You should now have a load balanced Credhub cluster you can interact with.

A quick health check can be run using:
```bash
curl -k https://<YOUR-LOADBALANCER-DNS-NAME>:8844/info
``` 

A healthy response should look like the following:

```json
{
  "auth-server": {
    "url": "https:your-uaa-url:8443"
  },
  "app": {
    "name": "CredHub",
    "version": "your-credhub-version"
  }
}
```

Congratulations. You are now ready to performance test Credhub.

## Running the Performance Test

Run performance tests via the bosh errand

```bash
bosh -d credhub-ha-perf run-errand credhub_cannon --download-logs
```

The errand logs contain the data needed to build the plot as a `csv` file. To build headroom plot, run:
 
```bash
python src/headroomplot/headroomplot.py <YOUR-DATA-FILE>.csv
```
> **Note**
> To install the required libraries for headroomplot, run `pip install -r src/headroomplot/requirements.txt`
