# CATS-DB


<img width="1080" height="360" alt="image" src="https://github.com/user-attachments/assets/a0eee31f-a595-4ece-85a0-15f431bc6891" />







A simple time series database built in Go on top of an underlying LSM storage engine, supporting: 

1. High throughput telemetry data ingestion.
2. Batched durability with WAL
3. MVCC through a lock free skip list.
4. Leveled compaction for key partiotioned SSTs.
5. Range read queries.
