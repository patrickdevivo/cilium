.. only:: not (epub or latex or html)

    WARNING: You are looking at unreleased Cilium documentation.
    Please use the official rendered version released here:
    http://docs.cilium.io

.. _k8s_install_aks:

*************************
Installation on Azure AKS
*************************

This guide covers installing Cilium into an Azure AKS environment. This guide
will work when setting up AKS in both `Basic <https://docs.microsoft.com/en-us/azure/aks/concepts-network#kubenet-basic-networking>`_ and `Advanced 
<https://docs.microsoft.com/en-us/azure/aks/concepts-network#azure-cni-advanced-networking>`_ networking mode.

This is achieved using Cilium in CNI chaining mode, with the Azure CNI plugin
as the base CNI plugin and Cilium chaining on top to provide L3-L7
observability, network policy enforcement enforcement, Kubernetes services
implementation, as well as other advanced features like transparent encryption
and clustermesh.

Prerequisites
=============

Ensure that you have the `Azure Cloud CLI 
<https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest>`_ installed.

To verify, confirm that the following command displays the set of available
Kubernetes versions.

.. code:: bash

        az aks get-versions -l westus -o table

Create an AKS Cluster
=====================

The full background on creating AKS clusters in advanced networking mode, see
`this guide <https://docs.microsoft.com/en-us/azure/aks/configure-azure-cni>`_.
You can use any method to create and deploy an AKS cluster with the exception
of specifying the Network Policy option. Doing so will still work but will
result in unwanted iptables rules being installed on all of your nodes.

If you want to us the CLI to create a dedicated set of Azure resources
(resource groups, networks, etc.) specifically for this tutorial, the following
commands (borrowed from the AKS documentation) run as a script or manually all
in the same terminal are sufficient.

It can take 10+ minutes for the final command to be complete indicating that
the cluster is ready.

.. note:: **Do NOT specify the '--network-policy' flag** when creating the cluster,
    as this will cause the Azure CNI plugin to push down unwanted iptables rules:

.. code:: bash

        export SP_PASSWORD=mySecurePassword
        export RESOURCE_GROUP_NAME=myResourceGroup-NP
        export CLUSTER_NAME=myAKSCluster
        export LOCATION=westus

        # Create a resource group
        az group create --name $RESOURCE_GROUP_NAME --location $LOCATION

        # Create a virtual network and subnet
        az network vnet create \
            --resource-group $RESOURCE_GROUP_NAME \
            --name myVnet \
            --address-prefixes 10.0.0.0/8 \
            --subnet-name myAKSSubnet \
            --subnet-prefix 10.240.0.0/16

        # Create a service principal and read in the application ID
        SP_ID_PASSWORD=$(az ad sp create-for-rbac --skip-assignment --query [appId,password] -o tsv)
        SP_ID=$(echo ${SP_ID_PASSWORD} | sed -e 's/ .*//g')
        SP_PASSWORD=$(echo ${SP_ID_PASSWORD} | sed -e 's/.* //g')
        unset SP_ID_PASSWORD

        # Wait 15 seconds to make sure that service principal has propagated
        echo "Waiting for service principal to propagate..."
        sleep 15

        # Get the virtual network resource ID
        VNET_ID=$(az network vnet show --resource-group $RESOURCE_GROUP_NAME --name myVnet --query id -o tsv)

        # Assign the service principal Contributor permissions to the virtual network resource
        az role assignment create --assignee $SP_ID --scope $VNET_ID --role Contributor

        # Get the virtual network subnet resource ID
        SUBNET_ID=$(az network vnet subnet show --resource-group $RESOURCE_GROUP_NAME --vnet-name myVnet --name myAKSSubnet --query id -o tsv)

        # Create the AKS cluster and specify the virtual network and service principal information
        # Enable network policy by using the `--network-policy` parameter
        az aks create \
            --resource-group $RESOURCE_GROUP_NAME \
            --name $CLUSTER_NAME \
            --node-count 1 \
            --generate-ssh-keys \
            --network-plugin azure \
            --service-cidr 10.0.0.0/16 \
            --dns-service-ip 10.0.0.10 \
            --docker-bridge-address 172.17.0.1/16 \
            --vnet-subnet-id $SUBNET_ID \
            --service-principal $SP_ID \
            --client-secret $SP_PASSWORD

Configure kubectl to Point to Newly Created Cluster
===================================================

Run the following commands to configure kubectl to connect to this
AKS cluster:

.. code:: bash

    az aks get-credentials --resource-group $RESOURCE_GROUP_NAME --name $CLUSTER_NAME


.. code:: bash

    export KUBECONFIG=/Users/danwent/.kube/config


To verify, you should see AKS in the name of the nodes when you run:

.. code:: bash

    kubectl get nodes
    NAME                       STATUS   ROLES   AGE     VERSION
    aks-nodepool1-12032939-0   Ready    agent   8m26s   v1.13.10

.. include:: k8s-install-azure-cni-steps.rst

.. include:: k8s-install-validate.rst
.. include:: hubble-install.rst
.. include:: getting-started-next-steps.rst

