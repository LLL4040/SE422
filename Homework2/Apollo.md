# Apollo

> Reading notes about Apollo: Scalable and Coordinated Scheduling for Cloud-Scale Computing 

## What is Apollo?

Apollo is a distributed scheduling framework for shared cluster state under Microsoft. Currently, Apollo has been deployed in Microsoft clustered production environment, is responsible for the scheduling and management of highly concurrent tasks million units on the number of server billions every day. The actual results, Apollo has good scalability, stability and good performance.

## What problem has Apollo solved?

* It  scale to make tens of thousands of scheduling decisions per second on a cluster with tens of thousands of servers; 
* It  maintain fair sharing of resources among different users and groups; 
* It  make high-quality scheduling decisions that take into account factors such as data locality, job characteristics, and server load, to minimize job latencies while utilizing the resources in a cluster fully. 

## Important methods used by Apollo

### 1. Distributed and (loosely) coordinated scheduling framework

​    To balance scalability and scheduling quality, Apollo adopts a distributed and (loosely) coordinated scheduling framework, in which independent scheduling decisions are made in an optimistic and coordinated manner by incorporating synchronized cluster utilization information.  Such a design strikes the right balance: it avoids the suboptimal (and often conflicting) decisions by independent schedulers of a completely decentralized architecture, while removing the scalability bottleneck and single point of failure of a centralized design. 

### 2. Apollo minimizes the task completion time 

​    To  achieve  high-quality  scheduling  decisions, Apollo schedules each task on a server that minimizes the task completion time.  The estimation model incorporates a variety of factors and allows a scheduler to perform a weighted decision,rather than solely considering data locality or server load. The data parallel nature of computation allows Apollo to refine the estimates of task execution time continuously based on observed runtime statistics from similar tasks during job execution. 

### 3. Each scheduler has information about the entire cluster for scheduling decisions 

​    Apollo introduces a lightweight hardware-independent  mechanism  to  advertise  load  on servers. When combined with a local task queue on each server, the mechanism provides a near-future view of resource availability on all the servers,which is used by the schedulers in decision making. Each scheduler has information about the entire cluster for scheduling decisions.

### 4. a series of correction mechanisms 

​    To cope with unexpected cluster dynamics, suboptimal estimations, and other abnormal runtime behaviors, which are facts of life in large-scale clusters,Apollo is made robust through a series of correction mechanisms that dynamically adjust and rectify suboptimal decisions at runtime.  

### 5. **opportunistic scheduling** 

​    Apollo divides operations into two categories, regular tasks and opportunistic tasks, ensuring low latency for regular jobs while using opportunistic jobs to increase cluster utilization and introducing a token-based mechanism to manage capacity. And to avoid overloading the cluster by limiting the total number of regular tasks.

## Apollo Framework

![apollo](https://github.com/PythonMyLife/SE100/blob/master/photos/apollo.png?raw=true)

​    Figure above provides an overview of Apollo’s architecture.  A Job Manager (JM), also called a scheduler, is assigned to manage the life cycle of each job. The global cluster load information used by each JM is provided through the cooperation of two additional entities in the Apollo framework: a Resource Monitor (RM) for each cluster and a Process Node (PN) on each server.  A PN process running on each server is responsible for managing the local resources on that server and performing local scheduling, while the RM aggregates load information from PNs across the cluster continuously, providing a global view of the cluster status for each JM to make informed scheduling decisions. 

​    The following figure depicts the communication between JM, PN and RM, showing how Apollo balances both distributed scheduling and shared cluster state.

![sequence](https://github.com/PythonMyLife/SE100/blob/master/photos/sequnce.png?raw=true)

## Performance of Apollo

### 1. Apollo at Scale

*  Apollo can constantly provide a scheduling rate of above 10,000, reaching up to 20,000 per second in a single cluster.  
*  Apollo is able to run 750 concurrent complex jobs (140,000 con-current regular tasks) and achieve over 90% CPU utilization when the demand is high during the weekdays,reaching closely the capacity of the cluster.  During the weekdays, 70% of Apollo’s workload comes from regular tasks. The balance shifts during the week-ends: more opportunistic tasks get executed on the avail-able resources when there are fewer regular tasks. 
* **Summary. ** Combined, those results show that Apollo is highly scalable, capable of scheduling over 20,000 re-quests per second, and driving high and balanced system utilization while incurring minimum syqueuing time. 

### 2. Scheduling Quality

*  Compare with the previously implemented baseline scheduler using production workload. And it shows that about 80% of recurring jobs receive various degrees of performance improvements. That's because  Apollo achieves much more balanced task queues across servers. 
*  Study business critical production jobs and use trace-based simulations to compare the quality. On average, the job latency improved around 22% with Apollo over the baseline scheduler, and Apollo performed within 4.5% of the oracle scheduler.
* **Summary**  Apollo delivers excellent job performance compared with the baseline scheduler and its scheduling quality is close to the optimal case. 

### 3. Evaluating Estimates 

*  Apollo provides good estimates on task wait time and CPU time, despite all the challenges, and estimation does help improve scheduling quality.

### 4. Correction Effectiveness 

*  Apollo’s duplicate scheduling is efficient, with 82% success rates, and accounts for less than 0.5% of task creations. And also Apollo is able to catch more than 70%stragglers efficiently and apply mitigation timely to expedite query execution.   Apollo’s correction mechanisms are shown effective with small overhead. 

### 5. Stable Matching Efficiency 

*  Apollo’s matching algorithm has the same asymptotic complexity as a naive greedy algorithm with negligible overhead. It performs significantly better than the greedy algorithm and is within 5% of the optimal scheduling in our simulation. 



## Personal thoughts

* Apollo is characterized by the fact that it uses a distributed architecture to achieve cluster state sharing, which avoids uneven server load, reduces scheduling conflicts, and achieves high scalability.
*  Shared state architecture may also cause some problems because state sharing is based on PN and RM constantly feedback operation information. In the case of high competition, the amount of information is so huge that JM needs extra resources to receive and process this information, which may affect the overall scheduling performance 

## Summary

​     Apollo is a cluster scheduling system that has been deployed in Microsoft's actual production environment and is responsible for the scheduling and management of millions of high-concurrency computing tasks on tens of thousands of machines every day. Apollo uses a policy based on shared state. Each scheduler has a global resource view with full access to the entire cluster. The Resource Monitor is responsible for maintaining a globally shared cluster state and providing each Job Manager with a global view of the cluster state.

## Reference

[1] [Apollo: Scalable and Coordinated Scheduling for Cloud-Scale Computing](https://www.usenix.org/system/files/conference/osdi14/osdi14-paper-boutin_0.pdf)

[2] [【每周论文】Apollo: Scalable and Coordinated Scheduling for Cloud-Scale Computing](https://blog.csdn.net/violet_echo_0908/article/details/78174782)

[3] [十大主流集群调度系统大盘点](https://blog.csdn.net/vip_iter/article/details/80123228)