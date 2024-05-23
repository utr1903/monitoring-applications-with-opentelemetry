# New Relic

## Getting started

New Relic is one of the leaders in the observability market. It's a highly scalable telemetry data platform tailored specifically for monitoring applications & infrastructures. You can go ahead and create completely free account [here](https://newrelic.com/signup). You'll have 1 Full User and 100GB free of ingest per each month.

After you login to your account, you can click on your avatar on the bottom left corner and click `API Keys`. After that, you can click on `Create a key` on the top right corner and create an `Ingest License Key`. With this key, you will be able send telemetry data to your New Relic account.

## Configuring the OTel collectors

The OpenTelemetry collector is capable of shipping the data (same or different) to multiple backends. This is achieved via the pipelines. Go ahead and comment in the `otlp/newrelic` related lines within the [configmap.yaml](/infra/helm/otelcollector/templates/configmap.yaml):

```yaml
otlp/newrelic:
  endpoint: https://otlp.nr-data.net:4317 # For EU datacenter -> https://otlp.eu01.nr-data.net:4317
  tls:
    insecure: false
  headers:
    api-key: <PASTE_YOUR_NEWRELIC_LICENSE_KEY>
```

Run the [`deploy.sh`](/infra/helm/deploy.sh) again for the otelcollectors and restart the pods:

```bash
kubectl delete pod -l app=otelcollector
```

The new instances of the collectors will be run with the new config and you will see the telemetry data flowing into your New Relic account in seconds!

## Entities

The New Relic world is built around a concept called _Entities_. The platform automatically creates you these entities and simultaneously builds the relationships in between them when the relevant telemetry data arrive to New Relic backend.

Let's take a look at our `httpserver` application for instance:

![Entity Summary](/media/newrelic_entitity_httpserver_summary.png)

This is a pre-built UI which gets generated automatically for you. On the left hand side, you have various capabilities that involves this entity. Let's take a look at the Distributed Tracing:

![Entity Distributed Tracing](/media/newrelic_entitity_httpserver_distributed_tracing.png)

The traces are automatically grouped for you. Since this is an HTTP server, they are grouped according to the REST endpoints. You can click on any one of the groups and then select one particular trace. You'd see the service map of that trace as well as the waterfall diagram of the spans:

![Entity Distributed Tracing individual](/media/newrelic_entitity_httpserver_distributed_tracing_individual.png)

## Dashboards

Although New Relic has a lot to offer out-of-the-box, you will want to create your own dashboards eventually where you want to put stuff according to whatever context you have in mind.

The cool thing is, New Relic provides you a programmability interface with which you can fully automate your New Relic resource creation. We will be using Terraform for that!

To talk to New Relic programmatically, we would need so called `User API Key`. You can click on your avatar on the bottom left corner and click `API Keys`. After that, you can click on `Create a key` on the top right corner and create an `User API Key`. With this key, you will be able query and provision things in your New Relic account.

Now let's refer to our Terraform scripts. Open up the [`deploy.sh`](terraform/deploy.sh) and put your account ID, your New Relic User API Key and your New Relic region:

```bash
# Initialize Terraform
terraform init

if [[ $flagDestroy != "true" ]]; then

  # Plan Terraform
  terraform plan \
    -var NEW_RELIC_ACCOUNT_ID=<YOUR_ACCOUNT_ID> \
    -var NEW_RELIC_API_KEY="<YOUR_API_KEY>" \
    -var NEW_RELIC_REGION="<YOUR_REGION>" \
    -out "./tfplan"

  # Apply Terraform
  if [[ $flagDryRun != "true" ]]; then
    terraform apply \
      -auto-approve \
      tfplan
  fi
else

  # Destroy Terraform
  terraform destroy \
    -auto-approve \
    -var NEW_RELIC_ACCOUNT_ID=<YOUR_ACCOUNT_ID> \
    -var NEW_RELIC_API_KEY="<YOUR_API_KEY>" \
    -var NEW_RELIC_REGION="<YOUR_REGION>"
fi
```

Let's deploy our dashboard to our account:

```bash
bash deploy.sh
```

### Performance

The Performance dashboard gives you an overview of how all of your applications are performing. It uses the metrics to generate the visualizations.

#### Throughput

![Throughput](/media/newrelic_dashboard_performance_throughput.png)

#### Latency

![Latency](/media/newrelic_dashboard_performance_latency.png)

#### Error rate

There will be no data at first because there are no errors happening yet.

### Troubleshooting

We will not be having a dashboard for this. In matter of fact, we will leverage the power of the New Relic Query Language (NRQL) to instantly find the root causes.

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

![Performance DB not reachable](/media/newrelic_dashboard_performance_error_rate_rpc.png)

Apparently, there is a problem with our GPRC applications and in matter of fact, all of the methods are failing. If we check the individual spans with errors, we see that the DB calls are all failing. It's a good start to check the spans that belong to DB calls.

Let's navigate to `All Entities` and pick the `grpcserver` from the Services tab. It will take us to the Summary page of the entity. Let's navigate to the Distributed Tracing UI. You will see that all of the trace groups started to have full errors at the same time. Let's pick on any one of the groups and investigate one of their traces:

![Troubleshooting DB not reachable](/media/newrelic_ui_troubleshooting_db_not_reachable_spans.png)

We see that our `grpcserver` application fails to talk to our database. Moreover, under the Error Details, you can see the Span Events. When you click on it, you'll see the `exception.message: Creating task failed. Database is not reachable.`.

Another cool feature is that the logs belonging to that trace are already correlated and put in context for you. Click the Logs tab on the top left corner and see the logs:

![Troubleshooting DB not reachable Trace](/media/newrelic_ui_troubleshooting_db_not_reachable_logs.png)

We can get the same information from the logs as well.

So, we did find the root cause but could we be faster? Yes indeed! Let's open up a Query Builder and replicate our debugging procedure above with NRQL.

First, execute the following query to get the spans with errors:

```
FROM Span SELECT * WHERE service.name = 'grpcserver' AND otel.status_code = 'ERROR'
```

You can use the `uniques()` function to get a hashset of values of an attribute. Do it for the `trace.id` to get all the trace IDs with errors.

Then, refer to the SpanEvent table and join it using the unique trace IDs:

```
FROM SpanEvent SELECT uniques(exception.message) WHERE trace.id IN (FROM Span SELECT uniques(trace.id) WHERE service.name = 'grpcserver' AND otel.status_code = 'ERROR')
```

You will get the following exception messages immediately:

```
Listing tasks failed. Database is not reachable.
Deleting tasks failed. Database is not reachable.
Creating task failed. Database is not reachable.
```

### Postprocessing failed

Now, we cause another issue! You can reset `grpcserver` btw. This time it'll be a postprocessing failure in the `httpclient` application. To do that, let's do the following change in the [values.yaml](../../infra/helm/httpclient/values.yaml):

```yaml
# Create postprocessing error
createPostprocessingError: true
```

Now, we re-run the [`deploy.sh`](/infra/helm/deploy.sh) script for the `httpclient` and wait for the application to be re-deployed.

Checking out our Performance dashboard, we can notice that the Error Rate widget regarding spans for the `httpclient` is now getting populated, though the HTTP ones don't.

![Performance postprocessing failed](/media/newrelic_dashboard_performance_error_rate_http.png)

This would mean that we the calls to the `httpserver` are successful, yet the something's happening in the postprocessing phase.

Let's navigate to `All Entities` and pick the `httpclient` from the Services tab. It will take us to the Summary page of the entity. Let's navigate to the Distributed Tracing UI. You will see that all of the trace groups started to have full errors at the same time. Let's pick on any one of the groups and investigate one of their traces:

![Troubleshooting postprocessing failed reachable](/media/newrelic_ui_troubleshooting_postprocessing_failed_spans.png)

Now we navigate to the Logs tab and figure out that the postprocessing calculation throws an exception during the calculation due to a singularity!

![Troubleshooting postprocessing failed reachable](/media/newrelic_ui_troubleshooting_postprocessing_failed_logs.png)

Let's do the same with NRQL:

```
FROM Log SELECT uniques(message) WHERE level = 'error' AND trace.id IN (FROM Span SELECT uniques(trace.id) WHERE service.name = 'httpclient' AND otel.status_code = 'ERROR')
```

### Postprocessing delayed

Time for the last issue! You can reset `httpclient`. This time it'll be a postprocessing delay in the `grpcclient` application. To do that, let's do the following change in the [values.yaml](../../infra/helm/grpcclient/values.yaml):

```yaml
# Create postprocessing delay
createPostprocessingDelay: true
```

This one is actually a bit tricky. When we check the Performance dashboard, the only thing which has changed seems to be the throughput.

![Performance postprocessing delayed](/media/newrelic_dashboard_performance_postprocessing.png)

It looks like it's almost halved but that does not mean that there is a problem. Maybe the traffic on our apps simply is reduced. Let's take a look at the logs:

```
FROM Log SELECT *
```

![Troubleshooting postprocessing delayed Logs](/media/newrelic_nrql_troubleshooting_postprocessing_delayed_logs.png)

We have a ton of warning logs coming from the `grpcclient` which state `Postprocessing schema cache could not be found. Calculating from scratch.`. Seems like there is an unexpected additional calculation being made.

Since we are enriching our logs with the trace context, we can grab the trace ID from one of these logs and explore that in Trace UI:

![Troubleshooting postprocessing delayed Trace](/media/newrelic_ui_troubleshooting_postprocessing_delayed_trace.png)

We can now see that the postprocessing takes about a second whereas the rest of the spans take no more than 250 milliseconds. Actually, this also explains why the throughput has decreased. When the amount of users using our apps remain the same but the latency increases, the request per minute they are able to make decreases with it.

Let's retrieve the exact same information only with NRQL:

```
FROM Log SELECT uniques(error.message) WHERE level !='info' AND trace.id IN (FROM Span SELECT uniques(trace.id) WHERE service.name = 'grpcclient' AND duration.ms > (FROM Span SELECT percentile(duration.ms, 90) WHERE service.name = 'grpcclient' AND span.kind = 'client'))
```

In the most inner query, we're calculating the 90th percentile of all of the client requests of the `grpcclient` application. Every span which takes longer than this worths an investigation.

So in the next outer query, we are filtering all of the client spans of the `grpcclient` application and getting the trace IDs of these requests.

Last, we gotta get the logs of these traces and for that we're having the most outer query against the Log table with joining the trace IDs from the spans.

Looks like we might want to check that cache out to improve the process!
