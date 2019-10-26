# Borg

> Reading notes about**《Large-scale cluster management at Google with Borg》**

## What is Borg?

Google’s Borg system is a cluster manager that runs hundreds of thousands of jobs, from many thousands of different applications, across a number of clusters each with up to tens of thousands of machines. 

Google uses Borg to admit, schedule, start, restart, and monitor the full range of applications that Google runs.

![Figure1](https://static.oschina.net/uploads/img/201510/10201735_pHVm.jpg)

## What problems has Borg solved?

* It hides the details of resource management and failure handling so its users can focus on application development instead.
* It operates with very high reliability and availability, and supports applications that do the same.
* It lets us run workloads across tens of thousands of machines effectively.

## How did Borg solve those problems?

### 1. Isolation

1. For users

   The machines in a cell are heterogeneous in many dimensions: sizes (CPU, RAM, disk, network), processor type, performance, and capabilities such as an external IP address or flash storage. Borg isolates users from most of these differences by determining where in a cell to run tasks, allocating their resources, installing their programs and other dependencies, monitoring their health, and restarting them if they fail. 

2. For tasks

   Google uses a Linux chroot jail as the primary security isolation mechanism between multiple tasks on the same machine.

   All Borg tasks run inside a Linux cgroup-based resource container and the Borglet manipulates the container settings, giving much improved control because the OS kernel is in the loop. The Borglet dynamically adjusts the resource caps of greedy LS tasks in order to ensure that they do not starve batch tasks for multiple minutes, selectively applying CFS bandwidth control when needed; shares are insufficient because there are multiple priority levels. 

### 2. Priority and quota

​    Every job has a *priority*, a small positive integer. A high priority task can obtain resources at the expense of a lower priority one, even if that involves preempting (killing) the latter. Borg defines non-overlapping *priority bands* for different uses, including (in decreasing-priority order): monitoring, production, batch, and best effort (also known as testing or free). 

​    In order to eliminate preemption cascades, they disallow tasks in the production priority band to preempt one another. And for overbuy, they respond to this by over-selling quota at lower-priority levels: every user has infinite quota at priority zero, although this is frequently hard to exercise because resources are oversubscribed. A low-priority job may be admitted but remain pending (unscheduled) due to insufficient resources. 

### 3.  Monitoring

​    Almost every task run under Borg contains a built-in HTTP server that publishes information about the health of the task and thousands of performance metrics (e.g., RPC latencies). Borg monitors the health-check URL and restarts tasks that do not respond promptly or return an HTTP error code. Other data is tracked by monitoring tools for dashboards and alerts on service level objective (SLO) violations. 

​    Borg records all job submissions and task events, as well as detailed per-task resource usage information in Infrastore, a scalable read-only data store with an interactive SQL-like interface via Dremel. This data is used for usage-based charging, debugging job and system failures, and long-term capacity planning. It also provided the data for the Google cluster workload trace.

### 4. Reliability

​    The Borgmaster is logically a single process but is actually replicated five times. Each replica maintains an in memory copy of most of the state of the cell, and this state is also recorded in a highly-available, distributed, Paxos-based store on the replicas’ local disks. A master is elected (using Paxos) when the cell is brought up and whenever the elected master fails. 

​    And *checkpoint* will also be stored in the Paxos store. Checkpoints can be used to restore a Borgmaster's state while debugging or for offline simulations. 

​    The Borgmaster polls each Borglet every few seconds to retrieve the machine’s current state and send it any outstanding requests. This gives Borgmaster control over the rate of communication, avoids the need for an explicit flow control mechanism, and prevents recovery storms. 

### 5. Schedule

​    The scheduling algorithm has two parts: *feasibility checking*, to find machines on which the task could run, and *scoring*, which picks one of the feasible machines. 

​    Google's current scoring model is a hybrid one that tries to reduce the amount of stranded resources – ones that cannot be used because another resource on the machine is fully allocated. It provides about 3–5% better packing efficiency than best fit for their workloads. 

​    The Borg scheduler uses limits to calculate feasibility for prod tasks. Score caching, equivalence classes and relaxed randomization make the Borg scheduler more scalable.

### 6. Availability

​    A key design feature in Borg is that already-running tasks continue to run even if the Borgmaster or a task’s Borglet goes down. 

​    Borg also does many things to mitigate the impact of failures, for example:

* automatically reschedules evicted tasks, on a new machine if necessary.
* reduces correlated failures by spreading tasks of a job across failure domains such as machines, racks, and power domains.
* limits the allowed rate of task disruptions and the number of tasks from a job that can be simultaneously down during maintenance activities such as OS or machine upgrades.
* uses declarative desired-state representations and idempotent mutating operations, so that a failed client can harmlessly resubmit any forgotten requests.
* rate-limits finding new places for tasks from machines that become unreachable, because it cannot distinguish between large-scale machine failure and a network partition.
* avoids repeating task::machine pairings that cause task or machine crashes.
* recovers critical intermediate data written to local disk by repeatedly rerunning a logsaver task, even if the alloc it was attached to is terminated or moved to another machine. Users can set how long the system keeps trying, a few days is common. 

### 7. Utilization

1. Cell sharing

   Most Borg cells are shared by thousands of users. From the experiments, even assuming the least-favorable of their results, sharing is still a win: the CPU slowdown is outweighed by the decrease in machines required over several different partitioning schemes, and the sharing advantages apply to all resources including memory and disk, not just CPU.

2. Large cells

   Google builds large cells, both to allow large computations to be run, and to decrease resource fragmentation.

3. Resource reclamation

   Rather than waste allocated resources that are not currently being consumed, Google estimates how many resources a task will use and reclaim the rest for work that can tolerate lower-quality resources, such as batch jobs. A machine may run out of resources at runtime if the reservations (predictions) are wrong – even if all tasks use less than their limits. If this happens, we kill or throttle non-prod tasks, never prod ones. 

## Summary

​    On the one hand, a number of Borg’s design features have been remarkably beneficial and have stood the test of time. On the other hand, there are also some shortcomings in Borg's design. Google recognized that and made alternative designs in Kubernetes. In my opinion, it seems, Borg, the cluster manager provided by Google, is too large, and it may have a long way to go for a really large range of use.

## Reference

* **《Large-scale cluster management at Google with Borg》** Abhishek Verma, Luis Pedrosa, Madhukar Korupolu, David Oppenheimer, Eric Tune, John Wilkes, Google Inc. 
* [在Google使用Borg进行大规模集群的管理 1-2](https://my.oschina.net/HardySimpson/blog/515398)
* [在Google使用Borg进行大规模集群的管理 3-4](https://my.oschina.net/HardySimpson/blog/516023)
* [在Google使用Borg进行大规模集群的管理 5-6](https://my.oschina.net/HardySimpson/blog/517283)
* [在Google使用Borg进行大规模集群的管理 7-8](https://my.oschina.net/HardySimpson/blog/518140)