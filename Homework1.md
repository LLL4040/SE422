# Technical report about Storage, Network, xPU and Memory

## Contents

* [Storage(Ceph)](#Storage(Ceph))
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

#### 1. Advantage

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

#### 2. Disadvantage

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

#### 1. Hardware planning

* Processor
    * The ceph-osd process consumes CPU resources during the running process, so it is usually bound to one CPU core for each ceph-osd process.Of course, if you use the EC method, you may need more CPU resources.
    * The ceph-mon process does not consume CPU resources very much, so you do not have to reserve too much CPU resources for the ceph-mon process.
    * Ceph-msd is also very CPU-intensive, so you need to provide more CPU resources.
* RAM
    * Ceph-mon and ceph-mds require 2G of memory, and each ceph-osd process requires 1G of memory, of course 2G is better.
* Network Planning
    * The 10G network is basically necessary to run Ceph. In network planning, we also try to separate the cilent and cluster networks.

#### 2. SSD

* In terms of cost, SATA SSD is generally chosen as the Journal.
* If the PCIE SSD is recommended when the budget is sufficient, performance will be further improved.

#### 3. Set the appropriate ceph configuration parameters

## Network
### 标题1
### 标题2
## xPU 
### fyr
### fyr

## Memory

### Supplier

#### 1. Memory chip manufacturers

* Korean: **Samsung**, **Hynix** (modern)
* Europe and the United States: **Micron**, ST, NXP, **Intel**, Infineon, SANDISK, TI (TI), ADI (Adejo), AMD (Super), SKYWORKS(Sica), Spansion ( Fesso Semiconductor), **Kingston** 
* Japanese: Elpida (acquired by Micron), **Toshiba**, Renesas,  **Panasonic**
* Taiwan: Huaya Technology (Meiguang Joint Venture), Nanya Technology, Ruijing (Elpida Joint Venture), Winbond, MTK (MediaTek), MXIC (Wang Hong), Etron (Zhuangchuang), Zentel (Li Ji, Force) Crystal Investment)

#### 2. Memory particle manufacturers

* Korean: Samsung, Hynix
* Japanese: NEC (Enyi), Hitachi (Hitachi), ELPIDA (Elpida),
* Europe and the United States: Micron, **Siemens**, Infineon   (Infineon, Germany, reorganized by Siemens memory), Qimonda (Qi  Mengda, Germany, reorganized by Infineon memory)
* Taiwan: Nanya (South Asia), Winbond (Winbond), PSC (Power Crystal), ProMos (Mao De)

#### 3. Terminal storage module and product manufacturer

* Qinmao (Shanghai), Kingmax (Sheng Chuang, Taiwan), Corsair (Jolly Roger) - DDR4, Apacer (Apacer), Jinbang, ADATA (ADATA)

#### 4. NAND Flash chip manufacturers

* **Samsung**, **Toshiba**, **Micron** (acquired Elpida), SK Hynix, SanDisk

#### 5. DRAM chip manufacturers

* **Samsung**, SK **Hynix**, **Micron**

#### 6.eMMC control IC manufacturers

* Huirong, group association, qingtai, Xinchuang 3S

### Technology type

#### 1.MCP (Multi-Chip Packaging; MCP)

Multi-chip packaging technology is to package two or more kinds of memory chips in the same BGA through integration (horizontal placement) and/or stacking. The storage scheme is generally **one NAND Flash chip** plus one low power consumption. **DRAM or Mobile RAM**. Its main application areas are handheld smart terminal devices such as **mobile phones**.

##### Features

Because it is packaged together, it can effectively reduce external interference, enhance the communication capability between the NAND Flash memory chip and DEAM, and improve the overall performance of the chip. Small size and low cost. In addition, the MCP package can save 70% space compared to 2 TSOPs, saving an average of 30%-40% board space for the final product. The MCP product not only saves board space, but also simplifies the manufacturing process and saves the original. It shortens the development time of new terminal products and speeds up the launch of the terminal products.

##### Application areas

Smartphone, feature phone, etc.

#### 2.eMMC

It is **a package of Flash control IC and NAND Flash**. This solution is better than NAND flash alone and has better compatibility. eMMC's embedded memory standard specification for the JEDEC Association has the advantage of simplifying memory design. The package uses a multi-chip package (MCP) technology to package the control IC and NAND Flash into a single component.

##### Advantage

eMMC combines data storage and replaces the NOR Flash support system boot function, which is widely used in smart phones, tablets, smart boxes, smart TVs and other products.

##### Application areas

**Tablets**, **smartphones**, etc.

#### 3.eMCP

It is based on MCP products. The storage solution is **an eMMC chip plus a low-power DRAM or Mobile RAM**.

##### Advantage

Compared with the traditional MCP, eMCP can reduce the burden of the main chip operation and manage the larger capacity of the flash memory because of **the built-in NAND Flash control chip**. In terms of external design, whether it is eMCP or eMMC The embedded memory design concept is to make the smart phone's exterior thinner and the cabinet more complete.

##### Application areas

Tablet, smartphone.

### Performance

#### 1.Speed

* The speed of the memory stick is generally used as the performance index when **the data is accessed once** (the unit is generally used in ns). The shorter the time, the faster the speed. Ordinary memory speed can only reach 70ns ~ 80ns, EDO memory speed can reach 60ns, and SDRAM memory speed has reached 7ns.
* Note: There are many manufacturers of memory modules. At present, there is no uniform labeling specification. Therefore, the memory performance indicators cannot be read out from the memory chip labels, but you can understand the speed, such as -70 or -60. The number indicates that the speed of this memory chip is 70ns or 60ns.

#### 2. Capacity

* There are various specifications for the size of the memory module. The early 30-line memory modules have 256K, 1M, 4M, and 8M capacities. The 72-line EDO memory is mostly 4M, 8M, and 16M, while the 168-line SDRAM memory is mostly 16M. , 32M, 64M, 128MB capacity, and even higher.

#### 3.CAS

* The CAS wait time refers to the period from when the read command is valid (issued on the rising edge of the clock) to when the output can provide data. It is usually 2 or 3 clock cycles, which determines the performance of the memory at the same operating frequency. The chip with a CAS latency of 2 is faster and performs better than a chip with a CAS latency of 3.

#### 4.GUF

* vThe higher frequency given by the manufacturer is lowered some, and the value thus obtained is called the rated available frequency GUF. For a memory bank of 8 ns, the highest usable frequency is 125 MHz, then the rated available frequency (GUF) should be 112 MHz. The highest available frequency and the rated available frequency (front-end system bus operating frequency) maintain a certain margin to maximize the stability of the system.

### How to choose

And when we choose, how big is it? One of the most important decision parameters is the maximum amount of memory supported by the CPU and the motherboard. Some CPUs and motherboards support 32G of memory, and more than 32G of memory is wasted, which is superfluous.

First, let's look at the speed of our own memory stick. If it is DDR4 or 2400MHz, then the memory stick you add is best at the same rate. A memory stick below 2400MHz will pull down the overall speed, even if you add a 16G size, it will not help. Similarly, buy a memory stick higher than 2400MHz, the operating state will be pulled down to 2400MHz.

The first is the maximum capacity that the CPU and the motherboard can support; the second is the reasonable matching of the card slot, the performance of the whole machine should be balanced; the last is to see the parameters of the computer memory, the same type of frequency is the best, too High and low will have an impact on the computer. Therefore, when we buy, don't be deceived by the merchants. The bigger the better, the only one that is most suitable for your computer is the one you choose.

## Reference

* [Computer data storage](https://en.wikipedia.org/wiki/Computer_data_storage#Further_reading)
* [Ceph (software)](https://en.wikipedia.org/wiki/Ceph_(software))
* [001 Ceph简介](https://www.cnblogs.com/zyxnhr/p/10530140.html)
* [关于 Ceph 现状与未来的一些思考](https://www.infoq.cn/article/some-thinking-about-the-present-situation-and-future-of-ceph)
* [Ceph 目前的优缺点分析](https://eof.net/p/5c4da57fe1382301b7843a4a)
* [Ceph的稳定性与性能调优](http://gcs.truan.wang/portfolio/ceph-stability-and-performance-tuning/)
* [Ceph性能优化总结(v0.94)](http://xiaoquqi.github.io/blog/2015/06/28/ceph-performance-optimization-summary/)
* [电脑内存条真的是越大越好吗？别再被商家骗钱了！](https://baijiahao.baidu.com/s?id=1612726385807021316&wfr=spider&for=pc)
* [内存的主要性能和指标有哪些?](https://zhidao.baidu.com/question/564059652.html)
* [DRAM、NAND FLASH、NOR FLASH三大存储器分析](http://www.elecfans.com/d/659117.html)
* [Memory芯片的应用攻略](https://blog.csdn.net/u010794281/article/details/46045157)
