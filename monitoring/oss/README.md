# OSS

The OSS is the default monitoring solution that consists of:

- Prometheus to store the metrics,
- Tempo to store the traces,
- Loki to store the logs,
- Grafana to query and visualize all of the above.

You need to port-forward the Grafana Kubernetes service to your local machine to login to Grafana:

```bash
kubectl port-forward svc/grafana 3000
```

You can then use Grafana on your browser on `http://localhost:3000`.

## Datasources & dashboards

The Helm chart of Grafana ([`values.yaml`](infra/helm/grafana/values.yaml)) already has all of the datasources and dashboards setup.

After you successfully deploy Grafana (and the necessary telemetry storage backends), you will be able to see the dashboards (Metrics, Traces and Logs).