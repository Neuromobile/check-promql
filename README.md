check-promql is a sensu check plugin that queries an Prometheus server and return the output in a [sensu check format](https://docs.sensu.io/sensu-core/0.29/reference/checks/#sensu-check-specification). It support multi value check.

### Running

```
check-promql -host my-prometheus-server.com -port 9090 -critical 0.5 -warninig 0.3 -query 'container_cpu_usage_seconds_total{pod_name=~"my-pod"}'
```

### Execution options

| Options | Description | Default |
| ------- | ----------- | ------- |
| -host   | Prometheus server host | "" |
| -port   | Prometheus server port | 9090 |
| -auth-basic-user | Prometheus server basic user | "" |
| -auth-basic-password | Prometheus server basic password | "" |
| -ssl | Prometheus server SSL? | false |
| -query  | Pass the query in promQL format | "" |
| -critical | Pass the critical threshold that will be evaluated | 0.0 |
| -warning | Pass the warning threshold that will be evaluated | 0.0 |
| -lt     | Change whether value is less than check | false |

### Build

```
git clone https://github.com/neuromobile/check-promql.git
cd check-promql
go build .
```
