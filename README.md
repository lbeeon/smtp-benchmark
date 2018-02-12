# Smtp Benchmark Tool

### Install 
```
go get -u github.com/lbeeon/smtp-benchmark
```
### Parameter

|  option  |  default  |  require  |  Description |
|----------|-------------|------|------|
| --host/ -h |  n/a | True | Target MTA. |
| --workers/ -w |  1 | False | Numbers of workers. |
| --nums/ -n |  1 | False | Numbers of mails. |
| --seconds/ -s |  -1 | False | If option provided, nums would be ignore. (second) |
| --size/ -b |  1 | False | Set the mail body size. (kb) |
| --eml/ -e |   | False | Use custom eml. |

### Usage
```
> Test one mail to target MTA
./smtp-benchmark -h 127.0.0.1:25

> Test how long does it take to send 1000 mails with 1 worker
./smtp-benchmark -h 127.0.0.1:25 -n 1000

> Test how long does it take to send 1000 mails with 20 workers
./smtp-benchmark -h 127.0.0.1:25 -w 20 -n 1000

> Test MTA throughput during 60 seconds
./smtp-benchmark -h 127.0.0.1:25 -w 20 -t 60

> Test one mail to target MTA with custom eml file
./smtp-benchmark -h 127.0.0.1:25 -e ./sample.eml
```

