# module-ingest

Targets:
* 10,000+ packet/sec TCP ingestion

* Pluggable subscribers (disk, AI, cloud)

* Ring buffer + disk-backed queue with write-ahead logging

* Built-in auth and identity layer

### Project Structure & File Breakdown

`/cmd/ – Binary Entrypoints`

| File	                  | Description                                                                                                     |
|------------------------|-----------------------------------------------------------------------------------------------------------------|
| server/main.go	 | Main TCP server. Initializes config, starts listener, queue, and consumer manager.                              |
| queue-worker/main.go	  | Optional standalone consumer that reads from disk queue and processes entries (used in recovery or batch jobs). |
| recover-dump/main.go	  | Utility to replay or inspect a corrupted .wal or .queue file for recovery/debugging.                            |

<br />

### /pkg/ – Core Modules (Reusable, Testable)
`network/ – TCP Handling`

| File            | 	Description                                                                                              |
|-----------------|-----------------------------------------------------------------------------------------------------------|
| listener.go	    | Initializes and manages TCP listener. Handles incoming connections via epoll/kqueue (platform-dependent). |
| connection.go	  | Per-connection logic: reads, buffers, error handling.                                                     |
| socket_opts.go	 | Low-level TCP options (SO_REUSEADDR, SO_LINGER, TCP_NODELAY).                                             |

`parser/ – Binary Protocol/Message Parsing`

| File         | 	Description                                                                                 |
|--------------|----------------------------------------------------------------------------------------------| 
| parser.go	   | Parses incoming binary packet format into a Message struct. Includes framing, length checks. |
| checksum.go	 | CRC32 or SHA256 checksum verification for tamper detection.                                  |

`ringbuffer/ – In-Memory Queue`

| File     | 	Description                                                                                                 |
|----------|--------------------------------------------------------------------------------------------------------------|
| ring.go	 | Lock-free ring buffer queue for in-memory storage before disk flush. Optimized for speed and bounded memory. |

`diskqueue/ – Persistent Queue System`

| File	       | Description                                                                                          |
|-------------|------------------------------------------------------------------------------------------------------|
| disk.go	    | Disk-backed queue that writes to files in blocks.                                                    |
| wal.go	     | Write-Ahead Logging system: guarantees recovery after crash.                                         |
| rotator.go	 | Handles rotating disk queue segments for performance and storage limits.                             |
| meta.go	    | Stores offsets, pointers, and queue state metadata. Ensures queue resumes exactly where it left off. |

`auth/ – Identity & Authorization`

| File    | 	Description                                                                                                   |
|---------|----------------------------------------------------------------------------------------------------------------|
| auth.go | 	Provides lightweight identity model (users, modules, cores). Handles token validation, scoped access control. |

`consumer/ – Processing Backends`

| File              | 	Description                                                                      |
|-------------------|-----------------------------------------------------------------------------------|
| manager.go	       | Central dispatcher that routes messages to enabled consumers.                     |
| disk_writer.go	   | Writes raw messages to structured file format.                                    |
| ai_consumer.go    | 	Optional: routes data to on-board AI models or FPGA-connected inference systems. |
| cloud_uploader.go | 	Pushes batched messages to S3/GCS/etc. when online.                              |
| indexer.go	       | Updates real-time search or time-series index (e.g. local TSDB or LSM-based DB).  |

`telemetry/ – Logging & Metrics`

| File	      | Description                                                                                          |
|------------|------------------------------------------------------------------------------------------------------|
| logger.go	 | Simple structured logger. Supports stdout, file, or circular memory buffer.                          |
| metrics.go | 	Tracks ingest rate, queue length, flush latency, and failure rates. Exported via CLI or raw socket. |

`config/ – Runtime Configuration`

| File      | 	Description                                                                                          |
|-----------|-------------------------------------------------------------------------------------------------------|
| config.go | 	Loads config from environment or file. Exposes TCP port, queue size, disk paths, auth settings, etc. |

`utils/ – Utilities`

| File    | 	Description                                         |
|---------|------------------------------------------------------|
| pool.go | 	Buffer pool to minimize allocs per message.         |
| time.go | 	Timestamps, duration helpers, monotonic time logic. |

`/internal/ – Private App Bootstrap Code`

| Path                  | 	Description                                                                                                         |
|-----------------------|----------------------------------------------------------------------------------------------------------------------|
| bootstrap/startup.go	 | Wires everything together: config → server → disk → consumer manager. Contains the full application lifecycle hooks. |

`/data/ – Disk Queue Files`

| Files |	Description |
|----------------|-------------------------------------------------------------|
| .wal |	Append-only WAL for crash recovery. |
| .queue	 | Main disk queue binary file. |
| .meta	| Stores queue pointers and offsets. |

`/scripts/ – Benchmarks, Testing, and CLI Tools`

| File           | Description                                                 |
|----------------|-------------------------------------------------------------|
| bench.sh	      | Bash script to benchmark throughput, memory, disk I/O.      |
| flood_test.go	 | Go script that sends 10,000 packets/sec for stress testing. |

`/tests/ – Unit & Integration Tests`

| File              | Description                                            |
|-------------------|--------------------------------------------------------|
| ingest_test.go	   | Validates ingestion path from network → memory → disk. |
| recovery_test.go	 | Tests queue restart and WAL recovery after crash.      |

**luh Security**

Token-based identity tied to physical modules or sensor Cores.

Optional MAC address or hardware fingerprinting.

Message integrity check via CRC or hash.

Built-in replay attack protection via sequence numbers.

**🚀 Performance Goals**
* Component	Target
* TCP connections	10,000+ simultaneous clients
* Throughput	≥ 10k packets/sec
* Disk flush	≤ 1ms latency (via batch + WAL)
* RAM usage	Bounded, with ring buffer fallback
* Recovery time	≤ 5s on crash
