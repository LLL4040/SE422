# Technical report about Storage, Network, xPU and Memory

## Contents

* [Storage](#Storage(Ceph))
* [Network](#Network)
* [xPU](#xPU)
* [Memory](#Memory)
* [Reference](#Reference)

## Storage(Ceph)

### What is storage

Computer data storage, often called storage or memory, is a technology consisting of computer components and recording media that are used to retain digital data. It is a core function and fundamental component of computers.

Generally the fast volatile technologies (which lose data when off power) are referred to as "memory", while slower persistent technologies are referred to as "storage".

This section will focus on the latter, namely persistent storage technologies.

### What is Ceph

In computing, Ceph (pronounced /ˈsɛf/) is a free-software storage platform, implements object storage on a single distributed computer cluster, and provides interfaces for object-, block- and file-level storage.

### Ceph design

![A high-level overview of the Ceph's internal organization](https://upload.wikimedia.org/wikipedia/commons/thumb/3/3e/Ceph_components.svg/630px-Ceph_components.svg.png)

Ceph employs five distinct kinds of daemons:

* Cluster monitors (ceph-mon) that keep track of active and failed cluster nodes, cluster configuration, and information about data placement and global cluster state.
* Object storage devices (ceph-osd) that use a direct, journaled disk storage (named BlueStore, since the v12.x release) or store the content of files in a filesystem (preferably XFS, the storage is named Filestore)
* Metadata servers (ceph-mds) that cache and broker access to inodes and directories inside a CephFS filesystem.
* HTTP gateways (ceph-rgw) that expose the object storage layer as an interface compatible with Amazon S3 or OpenStack Swift APIs
* Managers (ceph-mgr) that perform cluster monitoring, bookkeeping, and maintenance tasks, and interface to external monitoring systems and management (e.g. balancer, dashboard, Prometheus, Zabbix plugin)

All of these are fully distributed, and may run on the same set of servers. Clients with different needs can directly interact with different subsets of them.

### Ceph features

1. Advantage
    * CRUSH algorithm
        * Ceph abandoned the traditional centralized storage metadata addressing scheme and switched to the CRUSH algorithm for data addressing.
        * Based on the consistency hash, CRUSH considers the isolation of the disaster recovery domain and implements copy placement rules for various types of workloads, such as cross-machine room and rack awareness.
        * Ceph will assign the CRUSH rule set to the storage pool.When the Ceph client stores or retrieves data from the storage pool, Ceph automatically recognizes the CRUSH rule set and the top-level bucket in the rule for storing and retrieving data.When Ceph processes the CRUSH rule, it identifies the primary OSD that contains a PG, which allows the client to connect directly to the primary OSD for data read and write.
    * High availability
        * The number of copies of data in Ceph can be defined by the administrator, and the CRUSH algorithm can be used to specify the physical storage location of the replica to separate the fault domain, tolerate multiple fault scenarios and automatically attempt parallel repair.
        * At the same time, it supports strong and consistent copies, and the copies can be stored in the host, rack, computer room, and data center, so it is safe and reliable.
        * Storage nodes can be self-managed and automatically repaired. No single point of failure, there is a strong fault tolerance.
    * High scalability
        * Ceph is different from swift, and all client read and write operations go through the proxy node. Once the cluster concurrency increases, the proxy node can easily become a single point bottleneck. Ceph itself does not have a master node, it is easier to expand, and in theory, its performance will increase linearly as the number of disks increases.
    * Rich in features
        * Ceph supports three call interfaces: object storage, block storage, and file system mount. Three ways can be used together.
        * From the most basic horizontal scaling, dynamic scaling, redundant disaster recovery, load balancing, etc. of the distributed system, to the very practical rolling upgrade, multi-storage pool, delay deletion, etc. in the production environment, to the CephFS cluster and snapshot on the tall, erasure code, cross-storage pool cache, etc., ceph has many features.
2. Disadvantage
    * The problem of capacity expansion caused by non-central distributed metadata management, capacity expansion, PG changes, low data migration efficiency, and pains in operation and maintenance.
    * For object storage, the massive file size problem is not optimized, even the performance of bluestore is limited, and the usual file merging method is not adopted.

### Ceph key indicator

* External read and write bandwidth
    * Suppose we want to provide block storage services for 1000 virtual machines. Each virtual machine mounts a system disk, and 30% of the virtual machines mount an additional data disk. We define the system disk read and write bandwidth of each virtual machine.80MB/s, the data disk has high speed / super high speed / ordinary three, the average is still planned at 80MB/s, the total external bandwidth required for our ceph storage is 1000 * 80MB + 300 * 80MB = 10400MB/s peak load,It takes about 100Gb of total network bandwidth. If the storage cluster and the computing cluster are separated, 100Gb network connection is needed between the storage cluster and the computing cluster. The price of such network equipment is also very high. The cheaper solution is to deploy on a rack.A certain proportion of storage and computing servers, the corresponding problem is more complex cluster configuration and maintenance.
* Storage capacity
    * 1000 system disks are fixed at 20GB per disk, and data disks are from 20GB to 20TB. Under scientific analysis, small-capacity data accounts for the vast majority, and is designed according to the total of 300TB data disks. Then we need a total capacity of 320TB.
* IOPS
    * The number of reads and writes per second is calculated by 100IOPS per system disk. The data disk ultra-high performance disk may need 10,000+ IOPS, the high performance disk is calculated by 100 IOPS, the ordinary disk is calculated by 10 IOPS, and our IOPS load range is100,000 IOPS + 3000 IOPS to 100,000 IOPS + 3 million IOPS. The ultra-high performance disk has very high requirements for IOPS. It is mainly used for database reading and writing. If the demand is large, it should be built separately.

### How to choose

1. Hardware planning
    * Processor
        * The ceph-osd process consumes CPU resources during the running process, so it is usually bound to one CPU core for each ceph-osd process.Of course, if you use the EC method, you may need more CPU resources.
        * The ceph-mon process does not consume CPU resources very much, so you do not have to reserve too much CPU resources for the ceph-mon process.
        * Ceph-msd is also very CPU-intensive, so you need to provide more CPU resources.
    * RAM
        * Ceph-mon and ceph-mds require 2G of memory, and each ceph-osd process requires 1G of memory, of course 2G is better.
    * Network Planning
        * The 10G network is basically necessary to run Ceph. In network planning, we also try to separate the cilent and cluster networks.
2. SSD
    * In terms of cost, SATA SSD is generally chosen as the Journal.
    * If the PCIE SSD is recommended when the budget is sufficient, performance will be further improved.
3. Set the appropriate ceph configuration parameters

## Network
### 标题1
### 标题2
## xPU
### 标题1
### 标题2
## Memory
### 标题1
### 标题2

## Reference

* [Computer data storage](https://en.wikipedia.org/wiki/Computer_data_storage#Further_reading)
* [Ceph (software)](https://en.wikipedia.org/wiki/Ceph_(software))
* [001 Ceph简介](https://www.cnblogs.com/zyxnhr/p/10530140.html)
* [关于 Ceph 现状与未来的一些思考](https://www.infoq.cn/article/some-thinking-about-the-present-situation-and-future-of-ceph)
* [Ceph 目前的优缺点分析](https://eof.net/p/5c4da57fe1382301b7843a4a)
* [Ceph的稳定性与性能调优](http://gcs.truan.wang/portfolio/ceph-stability-and-performance-tuning/)
* [Ceph性能优化总结(v0.94)](http://xiaoquqi.github.io/blog/2015/06/28/ceph-performance-optimization-summary/)
