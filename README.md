# apiserver-load-tester
abuses the k8s api for testing purposes

## Your whole goal is to do something like this

```
time="2023-02-06T19:38:41-06:00" level=info msg="cm: latency for update on nueve 1.979602 seconds"
time="2023-02-06T19:38:41-06:00" level=info msg="cm: latency for update on uno 1.981863 seconds"
I0206 19:38:42.071846 1149095 request.go:690] Waited for 1.972334475s due to client-side throttling, not priority and fairness, request: PUT:https://0.0.0.0:41489/api/v1/namespaces/default/configmaps/diez
time="2023-02-06T19:38:42-06:00" level=info msg="cm: latency for update on diez 1.980783 seconds"
time="2023-02-06T19:38:42-06:00" level=info msg="cm: latency for update on seis 1.977924 seconds"
time="2023-02-06T19:38:42-06:00" level=info msg="cm: latency for update on cuatro 1.980142 seconds"
```