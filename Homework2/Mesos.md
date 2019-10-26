# Mesos

## 一、What problem has the Mesos paper solved

### current problems
* no framework will be optimal for all applications
* Two common solutions for sharing a cluster today are either to statically partition the cluster and run one frame- work per partition, or to allocate a set of VMs to each framework.**Unfortunately**, these solutions achieve nei- ther high utilization nor efficient data sharing.  
* Many frameworks, such as Hadoop and Dryad, employ a fine-grained resource sharing model.these frameworks can achieve high utilization, as jobs can rapidly scale when new nodes become available. **Unfortunately**, because these frameworks are de- veloped independently, there is no way to perform fine- grained sharing across frameworks, making it difficult to share clusters and data efficiently between them.

### Mesos:a share commodity cluster
* We present Mesos, a platform for sharing commodity clusters between multiple diverse cluster computing frameworks, such as Hadoop and MPI. 
* In this paper, we propose Mesos, a thin resource sharing layer that enables fine-grained sharing across diverse cluster computing frameworks, by giving frameworks a common interface for accessing cluster resources.

## 二、How does the Mesos paper solve these problems?
### sharing model&resource offers
* Mesos is built around two design elements: a fine-grained sharing model at the level of tasks, and a distributed scheduling mechanism called resource offers that delegates scheduling decisions to the frameworks. 
###  how to build a scalable and efficient system
* delegating control over scheduling to the frameworks
* resource offer encapsulates a bundle of resources that a framework can allocate on a cluster node to run tasks
## 三、What are the important methods Mesos use?
### 1、main components of Mesos
![main components of Mesos](https://img-blog.csdn.net/20161213211754172?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMTM2NDYxMg==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)
* The master implements fine-grained sharing across frameworks using resource offers. 
* Each framework running on Mesos consists of two components: a scheduler that registers with the master to be offered resources, and an executor process that is launched on slave nodes to run the framework’s tasks.
* Figure 3 shows an example of how a framework gets scheduled to run tasks.
![example](https://img-blog.csdn.net/20161213211838188?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMTM2NDYxMg==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)

### 2、Resource Allocation
* Mesos delegates allocation decisions to a pluggable al- location module, so that organizations can tailor alloca- tion to their needs.
*  two allocation modules: one that performs fair sharing based on a generalization of max-min fairness for multiple resources and one that implements strict priorities.Similar policies are used in Hadoop and Dryad.
### 3、Isolation
* Mesos provides performance isolation between framework executors running on the same slave by leveraging existing OS isolation mechanisms. 
* **platform-dependent**:we support multiple isolation mechanisms through pluggable isolation modules.
* **OS container technologies** :specifically Linux Containers and Solaris Projects. These technologies can limit the CPU, memory, network bandwidth, and (in new Linux kernels) I/O usage of a process tree.
### 4、Making Resource Offers Scalable and Robust
* task scheduling in Mesos is a distributed pro- cess, it needs to be efficient and robust to failures.
Mesos includes three mechanisms to help with this goal. 
1. **filters**:because some frameworks will always reject cer- tain resources, Mesos lets them short-circuit the rejection process and avoid communication by providing filters to the master
2. **count resources offered**: because a framework may take time to respond to an offer, Mesos counts resources offered to a framework towards its allocation of the cluster.
3. **rescinds & re-offers**:if a framework has not responded to an offer for a sufficiently long time, Mesos rescinds the offer and re-offers the resources to other frameworks.
### 5、Fault Tolerance
* When the active master fails, the slaves and schedulers connect to the next elected master and repopulate its state.
* Aside from handling master failures, Mesos reports node failures and executor crashes to frameworks’ schedulers.
*  Mesos allows a framework to register multiple schedulers such that when one fails, another one is notified by the Mesos master to take over.

## 四、How's the Behavior of Mesos

* we find that Mesos performs very well when frameworks can scale up and down elastically, tasks durations are homogeneous, and frameworks prefer all nodes equally.
### 1、emulate a centralized scheduler
* When different frameworks pre- fer different nodes, we show that Mesos can emulate a centralized scheduler that performs fair sharing across frameworks
### 2、handle heterogeneous task durations
 * we show that Mesos can handle heterogeneous task durations without impacting the performance of frameworks with short tasks
 * In particular, we consider a workload where tasks that are either short and long, where the mean duration of the long tasks is significantly longer than the mean of the short tasks. 
* To further alleviate the impact of long tasks, Mesos can be extended slightly to allow allocation policies to reserve some resources on each node for short tasks.
### 3、improve overall cluster utilization
* We also discuss how frameworks are incentivized to improve their performance under Mesos, and argue that these incentives also improve overall cluster utilization.
* **Short tasks**: revocation or simply due to failures
* **Scale elastically**:The ability of a framework to use resources as soon as it acquires them.
* **Do not accept unknown resources**:Frameworks are incentivized not to accept resources that they cannot use
### 4、some limitations of Mesos’s distributed scheduling model
#### Fragmentation
* When tasks have heterogeneous resource demands, a distributed collection of frameworks may not be able to optimize bin packing as well as a centralized scheduler
#### Interdependent framework constraints
* It is possible to construct scenarios where, because of esoteric interdependencies between frameworks (e.g., certain tasks from two frameworks cannot be colocated), only a sin- gle global allocation of the cluster performs well.
#### Framework complexity
* Using resource offers may make framework scheduling more complex. 
## 四、How's the effection of Mesos
### Overhead
* The MPI job took on average 50.9s without Mesos and 51.8s with Mesos, while the Hadoop job took 160s without Mesos and 166s with Mesos. In both cases, the overhead of using Mesos was less than 4%.
### Data Locality through Delay Scheduling
![measurements](https://img-blog.csdn.net/20161213212510831?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMTM2NDYxMg==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)
* Figure 8 shows averaged measurements from the 16 Hadoop instances across three runs of each scenario. 
* As expected, job performance im- proves with data locality: jobs run 1.7x faster in the 5s delay scenario than with static partitioning.
### Spark Framework
* the benefit of running iterative jobs using the specialized Spark framework.
![Figure 9](https://img-blog.csdn.net/20161213212544874?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMTM2NDYxMg==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)

### Mesos Scalability
![Scalability](https://img-blog.csdn.net/20161213212613437?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMTM2NDYxMg==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)
* the overhead remains small (less than one second) even at 50,000 nodes. In particular, this overhead is much smaller than the average task and job lengths in data center workloads. 
### Failure Recovery
* the MTTR was between 4 and 8 seconds, with 95% confidence intervals of up to 3s on either side.
### Performance Isolation
* In particular, using Containers to isolate CPU usage between a MediaWiki web server (consisting of multiple Apache processes running PHP) and a “hog” application (consisting of 256 processes spinning in infinite loops) shows on average only a 30% increase in request latency for Apache versus a 550% increase when running without Containers.
## 五、Summary
* existing frameworks can effectively share resources using Mesos, that Mesos enables the development of specialized frameworks providing major performance gains, such as Spark, and that Mesos’s simple design allows the system to be fault tolerant and to scale to 50,000 nodes.
* Mesos has great scalability. It can add a cluster framework that provides more abundant resource requirements for resources according to the corresponding model, and choose which task to use which kind of resources to run, so that their development is not constrained by the language provided by the system.
* Each model has its own application and inapplicability. It has advantages and disadvantages. The universality and efficiency cannot be combined. The new management model of the corresponding system is a trade-off between versatility and efficiency.

## 六、Reference
* [Mesos: A Platform for Fine-Grained Resource Sharing in the Data Center](http://static.usenix.org/events/nsdi11/tech/full_papers/Hindman_new.pdf)
* [Mesos: 数据中心中细粒度资源共享的平台](https://blog.csdn.net/u011364612/article/details/53613511)
* [为什么需要mesos](https://www.cnblogs.com/fxjwind/archive/2013/03/27/2984953.html)