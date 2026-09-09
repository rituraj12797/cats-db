# CATS-DB


```
  ,-.       _,---._ __  / \
 /  )    .-'       `./ /   \
(  (   ,'            `/    /|
 \  `-"             \'\   / |
  `.              ,  \ \ /  |
   /`.          ,'-`----Y   |
  (            ;        |   '
  |  ,-.    ,-'         |  /
  |  | (   |        hjw | /
  )  |  \  `.___________|/
  `--'   `--'
```






A simple time series database built in Go on top of an underlying LSM storage engine, supporting: 

1. High throughput telemetry data ingestion.
2. Batched durability with WAL
3. MVCC through a lock free skip list.
4. Leveled compaction for key partiotioned SSTs.
5. Range read queries.



