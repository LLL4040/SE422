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

### Suppliers

* Seagate
    * Seagate is a professional hard drive manufacturer.Launched a full range of hard drives for multi-domain applications. At present, Seagate has three major product lines,Personal storage hard drive,Enterprise storage hard drive and Consumer Storage Hard Disks The above three series can meet the application needs of different levels and different fields. Seagate has a complete variety, rapid technological innovation, and good quality in the storage hard disk market.
* Western Digital
    * Western Digital offers a broad range of technologies and products, including storage systems for data center environments, storage platforms and data center drives; for mobility, terminal and computing environments for automotive, connected home, industrial and IOT, smartphones andTablets, monitored embedded mobile flash cards, and built-in hard drives for computing, enterprise, gaming, NAS and surveillance devices. WD designs environmentally friendly and energy-efficient GreenPower technology hard drives for WD desktops, enterprises, CE and external hard drives.
* HITACHI
    * HiTACHI's traditional strengths are on small mobile hard drives and laptop hard drives, but as its product line continues to evolve: Ulltrastar, a 3.5-inch drive for servers and workstations.Deskstar: 3.5-inch hard drive for high-performance personal computers.1.8'', 2.5'' mobile hard drive, 1.8" and 2.5" hard drive for mobile processing and portable devices. Hard drive series for specific applications, such as streaming consumer electronics applications.
* Contrast between them
![对比1](https://github.com/LLL4040/images/blob/master/SE422_storage1.png?raw=true)
![对比2](https://github.com/LLL4040/images/blob/master/SE422_storage2.png?raw=true)

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

### Introduction

A computer network is a digital telecommunications network which allows nodes to share resources. In computer networks, computing devices exchange data with each other using connections (data links) between nodes. These data links are established over cable media such as wires or optic cables, or wireless media such as Wi-Fi.

### Type

![LAN and WAN](https://upload.wikimedia.org/wikipedia/commons/thumb/6/6e/LAN_WAN_scheme.svg/1200px-LAN_WAN_scheme.svg.png)

#### 1. LAN(local area network)

* Features
    * The geographic coverage covered by the network is relatively small. Usually no more than tens of kilometers, even in a building or a room.
    * The transmission rate of information is relatively high, ranging from 1 Mbps to 10 Mbps, which has reached 100 Mbps.
    * The management and management rights of the network belong to a certain unit.
* Advantages(wireless LAN)
    * It saves the need for difficult, costly and time-consuming wiring construction, reduces the impact on the surrounding environment of the construction, and saves a lot of economic costs.
    * Installation is simple and fast. Generally, as long as one or more Access Point (AP) devices are installed, a local area network covering the entire building or area can be established. In this signal coverage area, any location can access the network.Extremely convenient to use.
* Disadvantages(wireless LAN)
    * The rate is slow, generally susceptible to interference, and power is limited.
    * Transmission medium vulnerability and insufficient WEP.
    * Still not completely out of the wired network.

#### 2. MAN(metropolitan area network)

The coverage is between WAN and LAN

#### 3. WAN(wide area network)

* Features
    * The WAN consists of more than two LANs, and the connections between these LANs can traverse a distance of more than 30 mile *.
    * A large WAN can consist of many LANs and MANs on all continents.
    * The most widely known WA N is the Internet, which is made up of thousands of LANs and WANs around the world.
* Advantages
    * Privacy and security: As mentioned above, a WAN provides a direct, dedicated connection through which your data can pass. This limits opportunities for others to intercept your data as it is in transit between locations. WANs offer a distinct privacy and security advantage.
    * Network performance consistency:  Equally important, your data does not have to compete with other Internet data for bandwidth as your communications travel between destinations. You get continuous access to all the bandwidth you are paying for. This allows your business to avoid throughput lags commonly experienced in the Internet.
    * Infrastructure management: The consistent performance and enhanced security noted above create a network environment in which you can effectively incorporate telecommunication solutions, such as VoIP. You therefore get the benefit of centralized management from a single dashboard for both your data network and telecommunications.
* Disadvantages
    * Initial budget costs and maintenance costs may be high.
    * WANs offer privacy advantages, but all networks have vulnerabilities.

### Key indicators

#### 1. Rate

* The rate, or data rate or bit rate, is one of the most important performance metrics in computer networks.
* The unit of rate is b/s, or kb/s, Mb/s, Gb/s, etc.The rate is often referred to as the nominal rate or nominal rate.

#### 2. Bandwidth

* In computing, bandwidth is the maximum rate of data transfer across a given path. Bandwidth may be characterized as network bandwidth, data bandwidth, or digital bandwidth.
* This table shows the maximum bandwidth (the physical layer net bitrate) of common Internet access technologies.
    | Technologies | Maximum bandwidth |
    | :----: | :----: |
    | 100 Gigabit Ethernet | 100 Gbit/s |
    | Thunderbolt 3 | 40 Gbit/s |
    | 10 Gigabit Ethernet, USB 3.1 | 10 Gbit/s |
    | OC192 | 9.6 Gbit/s |
    | Wireless 802.11ad | 7 Gbit/s |
    | USB 3.0 | 5 Gbit/s |
    | OC48 | 2.5 Gbit/s |
    | Wireless 802.11ac | 1.3 Gbit/s |
    | Gigabit Ethernet | 1 Gbit/s |
    | Fast Ethernet | 100 Mbit/s |
    | Ethernet | 10 Mbit/s |
    | Modem / Dialup | 56 kbit/s |
* Bandwidth is measured in bits per second.

#### 3. Throughput

* Throughput represents the amount of data that passes through a network (or channel, interface) per unit of time.
* Throughput is more often used to measure a network in the real world in order to know how much data is actually going through the network.
* The throughput is limited by the bandwidth of the network or the rated rate of the network.

#### 4. Delay

* When data is sent, the time it takes for the block to travel from the node to the transmission medium.That is, the time from the first bit of the transmitted data frame to the time when the last bit of the frame is transmitted.
* For high-speed network links, what we increase is only the rate at which data is sent, not the rate at which bits travel on the link.Increasing the link bandwidth reduces the transmission delay of the data.

#### 5. Delay bandwidth product

* In data communications, the bandwidth-delay product is the product of a data link's capacity (in bits per second) and its round-trip delay time (in seconds).
* The result, an amount of data measured in bits (or bytes), is equivalent to the maximum amount of data on the network circuit at any given time, i.e., data that has been transmitted but not yet acknowledged.
* A network with a large bandwidth-delay product is commonly known as a long fat network (shortened to LFN). As defined in RFC 1072, a network is considered an LFN if its bandwidth-delay product is significantly larger than 10^5 bits (12,500 bytes).

#### 6. Round trip time

* In telecommunications, the round-trip delay time (RTD) or round-trip time (RTT) is the length of time it takes for a signal to be sent plus the length of time it takes for an acknowledgement of that signal to be received. This time delay includes the propagation times for the paths between the two communication endpoints.
* In the context of computer networks, the signal is generally a data packet, and the RTT is also known as the ping time. An internet user can determine the RTT by using the ping command.

### How to choose

* If your organization consists of more than one person with a computer, and each needs to communicate with the others, then give serious consideration to a LAN.
* Choose according to fault tolerance, network management, performance, Internet connectivity, etc.

## xPU

### CPU

![cpu](https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1571566100&di=4a0ebd3c41ed12d537090f027755248a&imgtype=jpg&er=1&src=http%3A%2F%2Fimg1.cache.netease.com%2Fcatchpic%2F9%2F98%2F98FFADA8323A52C4BA2406D99A4639D5.jpg)

#### 1、CPU Vendor and technology type

* The current common computer CPU manufacturers are: Intel and AMD.
* Intel is the big brother of CPU production, it has more than 80% market share, Intel's CPU has become the de facto x86CPU technical specifications and standards. The most recent ones are the Celeron Athlon series, the Pentium Pentium series, the Core 2 Core series, and the i7 i5 i3 series CPU.
* AMD, the latest Athlon64x2 and Sempron have a good price/performance ratio, especially with 3DNOW+ technology, which makes it perform well on 3D.
* The following are rarely seen:IBM and Cyrix,IDT Corporation,IA VIA Technologies,Domestic dragon core.GodSon is a general-purpose processor with state-owned intellectual property rights. It has two generations of products, which can reach the low-end CPU level of INTEL and AMD on the market.
![intel and amd](https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1571566701&di=775da12589e52eda1a6ec78d12965f7a&imgtype=jpg&er=1&src=http%3A%2F%2Farticle.fd.zol-img.com.cn%2Ft_s640x2000%2Fg5%2FM00%2F00%2F06%2FChMkJloCzbGIW5b_AAKiFrfOCoQAAh8QgJQagkAAqIu206.jpg)

#### 2、CPU key indicator

* CPU cache capacity and performance: The larger the calculated cache capacity, the better his performance.
* CPU working voltage: the range of normal operating voltage of the CPU is wider.Now the voltage required for the normal operation of the CPU is getting lower and lower than before, and the current minimum is 1.1V.
* Strengthen the working voltage, strengthen the operating efficiency of the CPU, and achieve the purpose of overclocking.
* For the CPU, the indicators that affect its performance mainly include the main frequency, the number of CPU bits, and the CPU's cache instruction set.

#### 3、How to weigh and choose CPU and GPU

* The workflow and physical structure of the CPU and GPU are roughly similar. Compared to the CPU, the GPU works more singly.
* The CPU has a powerful arithmetic unit that can perform arithmetic calculations in a small number of clock cycles.The GPU is based on a large throughput design with many arithmetic units and very little cache.
* Because the CPU has a large number of caches and complex logic control units, it is very good at logic control, serial operations.
* Because GPU has a large number of arithmetic units, it can perform a lot of calculations at the same time. It is good at large-scale concurrent calculations. It has a large amount of calculation but no technical content, and it has to be repeated many times.
* to sum up, using the CPU to do complex logic control, using the GPU to do simple but large arithmetic operations, can greatly improve the speed of the program.

### NPU

![nup](https://gss0.baidu.com/-4o3dSag_xI4khGko9WTAnF6hhy/zhidao/wh%3D600%2C800/sign=7d49746fc611728b30788424f8cceff6/2e2eb9389b504fc2792b4bfbeedde71191ef6ddc.jpg)

#### 1、NPU Vendor and technology type

* The embedded neural network processor (NPU) adopts the "data-driven parallel computing" architecture, and is especially good at processing massive multimedia data of video and image.  
* Characteristics：Under the HiAI architecture, the AI ​​performance density is much better than that of the CPU and GPU, enabling more tasks to be completed faster with less power consumption, greatly improving the computational efficiency of the chip.
* Disadvantage：Loss of versatility in exchange for high performance

#### 2、NPU key indicator

* Performance and power consumption：This NPU evolution is similar to the previous architecture of CPUs from single core to large and small cores. The big core and small core are based on the same architecture design.  when a high-load task occurs, the large core can be topped up with strong performance, thereby enabling the chip to have strong performance while reducing the power consumption of the chip.
* compare with cpu and gpu（HUAWEI Kirin 970 as example）For the first time, the Kirin 970 integrates an NPU specifically tailored for deep learning, with FP16 performance of 1.92 TFLOP. Specifically, the NPU is 25 times that of the CPU and 6.25 times (25/4) of the GPU. The energy efficiency ratio is 50 times that of the CPU and 6.25 times (50/8) of the GPU.
![Kirin 970](https://gss1.bdstatic.com/9vo3dSag_xI4khGkpoWK1HF6hhy/baike/crop%3D49%2C0%2C561%2C371%3Bc0%3Dbaike80%2C5%2C5%2C80%2C26/sign=8da5a05aebcd7b89fd2360c332107b8b/aec379310a55b319b26622e54ea98226cefc1786.jpg)  

#### 3、How to weigh and choose NPU and others

* The NPU subverts the traditional von Neumann computer architecture used by the CPU.
* This type of data stream processor greatly increases the ratio of computing power to power consumption, enabling artificial intelligence to be used in embedded machine vision applications.
* while such specialized hardware can achieve greater efficiency in terms of hardware and power in handling artificial intelligence, it also loses flexibility.

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

* Multi-chip packaging technology is to package two or more kinds of memory chips in the same BGA through integration (horizontal placement) and/or stacking. The storage scheme is generally **one NAND Flash chip** plus one low power consumption. **DRAM or Mobile RAM**. Its main application areas are handheld smart terminal devices such as **mobile phones**.
* Features
    * Because it is packaged together, it can effectively reduce external interference, enhance the communication capability between the NAND Flash memory chip and DEAM, and improve the overall performance of the chip. Small size and low cost. In addition, the MCP package can save 70% space compared to 2 TSOPs, saving an average of 30%-40% board space for the final product. The MCP product not only saves board space, but also simplifies the manufacturing process and saves the original. It shortens the development time of new terminal products and speeds up the launch of the terminal products.
* Application areas
    * Smartphone, feature phone, etc.

#### 2.eMMC

* It is **a package of Flash control IC and NAND Flash**. This solution is better than NAND flash alone and has better compatibility. eMMC's embedded memory standard specification for the JEDEC Association has the advantage of simplifying memory design. The package uses a multi-chip package (MCP) technology to package the control IC and NAND Flash into a single component.
* Advantage
    * eMMC combines data storage and replaces the NOR Flash support system boot function, which is widely used in smart phones, tablets, smart boxes, smart TVs and other products.
* Application areas
    * **Tablets**, **smartphones**, etc.

#### 3.eMCP

* It is based on MCP products. The storage solution is **an eMMC chip plus a low-power DRAM or Mobile RAM**.
* Advantage
    * Compared with the traditional MCP, eMCP can reduce the burden of the main chip operation and manage the larger capacity of the flash memory because of **the built-in NAND Flash control chip**. In terms of external design, whether it is eMCP or eMMC The embedded memory design concept is to make the smart phone's exterior thinner and the cabinet more complete.
* Application areas
    * Tablet, smartphone.

### Performance

#### 1.Speed

* The speed of the memory stick is generally used as the performance index when **the data is accessed once** (the unit is generally used in ns). The shorter the time, the faster the speed. Ordinary memory speed can only reach 70ns ~ 80ns, EDO memory speed can reach 60ns, and SDRAM memory speed has reached 7ns.
* Note: There are many manufacturers of memory modules. At present, there is no uniform labeling specification. Therefore, the memory performance indicators cannot be read out from the memory chip labels, but you can understand the speed, such as -70 or -60. The number indicates that the speed of this memory chip is 70ns or 60ns.

#### 2. Capacity

* There are various specifications for the size of the memory module. The early 30-line memory modules have 256K, 1M, 4M, and 8M capacities. The 72-line EDO memory is mostly 4M, 8M, and 16M, while the 168-line SDRAM memory is mostly 16M. , 32M, 64M, 128MB capacity, and even higher.

#### 3.CAS

* The CAS wait time refers to the period from when the read command is valid (issued on the rising edge of the clock) to when the output can provide data. It is usually 2 or 3 clock cycles, which determines the performance of the memory at the same operating frequency. The chip with a CAS latency of 2 is faster and performs better than a chip with a CAS latency of 3.

#### 4.GUF

* The higher frequency given by the manufacturer is lowered some, and the value thus obtained is called the rated available frequency GUF. For a memory bank of 8 ns, the highest usable frequency is 125 MHz, then the rated available frequency (GUF) should be 112 MHz. The highest available frequency and the rated available frequency (front-end system bus operating frequency) maintain a certain margin to maximize the stability of the system.

### How to choose

* And when we choose, how big is it? One of the most important decision parameters is the maximum amount of memory supported by the CPU and the motherboard. Some CPUs and motherboards support 32G of memory, and more than 32G of memory is wasted, which is superfluous.
* First, let's look at the speed of our own memory stick. If it is DDR4 or 2400MHz, then the memory stick you add is best at the same rate. A memory stick below 2400MHz will pull down the overall speed, even if you add a 16G size, it will not help. Similarly, buy a memory stick higher than 2400MHz, the operating state will be pulled down to 2400MHz.
* The first is the maximum capacity that the CPU and the motherboard can support; the second is the reasonable matching of the card slot, the performance of the whole machine should be balanced; the last is to see the parameters of the computer memory, the same type of frequency is the best, too High and low will have an impact on the computer. Therefore, when we buy, don't be deceived by the merchants. The bigger the better, the only one that is most suitable for your computer is the one you choose.

## Reference

* [Computer data storage](https://en.wikipedia.org/wiki/Computer_data_storage#Further_reading)
* [Ceph (software)](https://en.wikipedia.org/wiki/Ceph_(software))
* [001 Ceph简介](https://www.cnblogs.com/zyxnhr/p/10530140.html)
* [关于 Ceph 现状与未来的一些思考](https://www.infoq.cn/article/some-thinking-about-the-present-situation-and-future-of-ceph)
* [Ceph 目前的优缺点分析](https://eof.net/p/5c4da57fe1382301b7843a4a)
* [Ceph的稳定性与性能调优](http://gcs.truan.wang/portfolio/ceph-stability-and-performance-tuning/)
* [Ceph性能优化总结(v0.94)](http://xiaoquqi.github.io/blog/2015/06/28/ceph-performance-optimization-summary/)
* [西部数据百度百科](https://baike.baidu.com/item/%E8%A5%BF%E9%83%A8%E6%95%B0%E6%8D%AE/6572761?fr=aladdin)
* [日立硬盘百度百科](https://baike.baidu.com/item/%E6%97%A5%E7%AB%8B%E7%A1%AC%E7%9B%98/6463486?fr=aladdin)
* [常见硬盘供应商及产品介绍](https://www.taodocs.com/p-288494283.html)
* [电脑内存条真的是越大越好吗？别再被商家骗钱了！](https://baijiahao.baidu.com/s?id=1612726385807021316&wfr=spider&for=pc)
* [内存的主要性能和指标有哪些?](https://zhidao.baidu.com/question/564059652.html)
* [DRAM、NAND FLASH、NOR FLASH三大存储器分析](http://www.elecfans.com/d/659117.html)
* [Memory芯片的应用攻略](https://blog.csdn.net/u010794281/article/details/46045157)
* [Computer network](https://en.wikipedia.org/wiki/Computer_network)
* [计算机网络基本概述及简单介绍](https://blog.csdn.net/qq_36474990/article/details/78743439)
* [网络](https://baike.baidu.com/item/%E7%BD%91%E7%BB%9C/143243#5)
* [无线局域网的优缺点](https://wenku.baidu.com/view/594d572e7375a417866f8f7a.html)
* [Wide Area Networks: The Pros and Cons](https://www.firstlight.net/wide-area-networks-the-pros-and-cons/)
* [Bandwidth (computing)](https://en.wikipedia.org/wiki/Bandwidth_(computing))
* [网络技术基础](https://www.jianshu.com/p/fc45e0ff73f4)
* [Bandwidth-delay product](https://en.wikipedia.org/wiki/Bandwidth-delay_product)
* [Round-trip delay time](https://en.wikipedia.org/wiki/Round-trip_delay_time)
* [Selecting the Right Computer Network](https://www.journalofaccountancy.com/issues/1997/feb/select.html)
* [Cpu introduction](https://baike.baidu.com/item/中央处理器/284033?fromtitle=CPU&fromid=120556&fr=aladdin)
* [Compare different cpu manufacturers](http://blog.sina.com.cn/s/blog_6e8cfd7b0102ysut.html)
* [Introduction to NPU and performance analysis of Kirin 970's NPU](http://www.elecfans.com/d/622258.html)
* [Why does NPU move toward heterogeneity?](http://www.wyzxwk.com/Article/chanye/2019/09/408023.html)
