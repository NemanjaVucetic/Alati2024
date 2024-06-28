# PromQL

- Upitni jezik za Prometheus
- Dokumentacija: https://prometheus.io/docs/prometheus/latest/querying/basics/

## Primeri

### Broj obradjenih HTTP zahteva u poslednjih 24

    increase(my_app_http_hit_total[24h])

### Broj uspesnih zahteva

    sum(increase(response_status_total{status=~"2..|3.."}[24h]))

### Broj nuuspesnih zahteva

    sum(rate(response_status_count{status=~"4..|5.."}[24h]))

### Prosecno vreme izvrsavanja zahteva za svaki endpoint

    histogram_quantile(0.95, sum(rate(http_response_time_bucket[24h])) by (le, endpoint))