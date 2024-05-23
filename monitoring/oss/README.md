# OSS

The OSS is the default monitoring solution that consists of:

- Prometheus to store the metrics,
- Tempo to store the traces,
- Loki to store the logs,
- Grafana to query and visualize all of the above.

You need to port-forward the Grafana Kubernetes service to your local machine to login to Grafana:

```bash
kubectl port-forward svc/grafana 3000:80
```

You can then use Grafana on your browser on `http://localhost:3000`.

## Datasources & dashboards

The Helm chart of Grafana ([`values.yaml`](infra/helm/grafana/values.yaml)) already has all of the datasources and dashboards setup.

After you successfully deploy Grafana (and the necessary telemetry storage backends), you will be able to see the dashboards (Performance and Troubleshooting).

### Performance

The Performance dashboard gives you an overview of how all of your applications are performing. It uses the metrics to generate the visualizations.

#### Throughput

![Throughout](/media/grafana_dashboard_performance_throughput.png)

#### Latency

![Latency](/media/grafana_dashboard_performance_latency.png)

#### Error rate

There will be no data at first because there are no errors happening yet.

### Troubleshooting

This dashboard stands for our debugging purposes where will be diving deep into our applications. Therefore, we will be refering to traces and logs.

In order to fasten our troubleshooting, there are a couple of dashboard variables built into the dashboard which we can set up to filter the necessary information we need. Moreover, we are only considering the logs with type `warning` and `error` in order to reduce the noise. Since, we have no error yet, there are no logs shown at first.

![Latency](/media/grafana_dashboard_troubleshooting.png)

## Generating issues & Troubleshooting

Now let's generate some intentional errors and try to debug them! Of course, we'll act as if we have no idea about the error.

### Database not reachable

Let's mock as if the `grpcserver` application cannot reach to the database. To do that, let's do the following change in the [values.yaml](../../infra/helm/grpcserver/values.yaml):

```yaml
# Create DB not reachable error
createDbNotReachableError: true
```

Now, we re-run the [`deploy.sh`](/infra/helm/deploy.sh) script for the `grpcserver` and wait for the application to be re-deployed.

Checking out our Performance dashboard, we can notice that the Error Rate widgets for the `grpcclient` and `grpcserver` are now being populated, indicating that there are errors.

![Latency](/media/grafana_dashboard_performance_error_rate_rpc.png)

Apparently, there is a problem with our GPRC applications and in matter of fact, all of the methods are failing. If we check the individual spans with errors, we see that the DB calls are all failing. It's a good start to check the spans that belong to DB calls.

We go to our Troubleshooting dashboard and apply the filters:

```
service: grpcserver
span_status: error
```

![Troubleshooting DB not reachable](/media/grafana_dashboard_troubleshooting_db_not_reachable.png)

We see that our `grpcserver` application fails to talk to our database. Moreover, we can go ahead and check one of the failed traces by clicking on the trace ID. It will pop up the Explore page with Tempo:

![Troubleshooting DB not reachable Trace](/media/grafana_explore_troubleshooting_db_not_reachable.png)

We can see that the error from the DB has propapagated up till the `grpcclient`. Plus, the error message is built into the span as a span event where the exception message also indicates what we've seen within the logs!

### Postprocessing failed

Now, we cause another issue! You can reset `grpcserver` btw. This time it'll be a postprocessing failure in the `httpclient` application. To do that, let's do the following change in the [values.yaml](../../infra/helm/httpclient/values.yaml):

```yaml
# Create postprocessing error
createPostprocessingError: true
```

Now, we re-run the [`deploy.sh`](/infra/helm/deploy.sh) script for the `httpclient` and wait for the application to be re-deployed.

Checking out our Performance dashboard, we can notice that the Error Rate widget regarding spans for the `httpclient` is now getting populated, though the HTTP ones don't.

![Performance postprocessing failed](/media/grafana_dashboard_performance_error_rate_http.png)

This would mean that we the calls to the `httpserver` are successful, yet the something's happening in the postprocessing phase.

Let's take a look at the Troubleshooting dashboard with the variables:

```
service: grpcserver
span_status: error
```

![Troubleshooting postprocessing failed](/media/grafana_dashboard_troubleshooting_postprocessing_failed.png)

Checking out the logs, we figure out that the postprocessing calculation throws an exception during the calculation due to a singularity!

Let's check one of the traces to see the waterfall diagram of the error:

![Troubleshooting postprocessing failed Trace](/media/grafana_explore_troubleshooting_postprocessing_failed.png)

### Postprocessing delayed

Time for the last issue! You can reset `httpclient`. This time it'll be a postprocessing delay in the `grpcclient` application. To do that, let's do the following change in the [values.yaml](../../infra/helm/grpcclient/values.yaml):

```yaml
# Create postprocessing delay
createPostprocessingDelay: true
```

This one is actually a bit tricky. When we check the Performance dashboard, the only thing which has changed seems to be the throughput.

![Performance postprocessing delayed](/media/grafana_dashboard_performance_throughput_postprocessing.png)

It looks like it's almost halved but that does not mean that there is a problem. Maybe the traffic on our apps simply is reduced. Let's take a look at the logs:

![Troubleshooting postprocessing delayed](/media/grafana_dashboard_troubleshooting_postprocessing_delayed.png)

We have a ton of warning logs coming from the `grpcclient` which state `Postprocessing schema cache could not be found. Calculating from scratch.`. Seems like there is an unexpected additional calculation being made. Since we are enriching our logs with the trace context, we can grab the trace ID from one of these logs and explore that in Trace UI:

![Troubleshooting postprocessing delayed Trace](/media/grafana_explore_troubleshooting_postprocessing_delayed.png)

We can now see that the postprocessing takes about a second whereas the rest of the spans take no more than 250 milliseconds. Actually, this also explains why the throughput has decreased. When the amount of users using our apps remain the same but the latency increases, the request per minute they are able to make decreases with it.

Looks like we might want to check that cache out to improve the process!
