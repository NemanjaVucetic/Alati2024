# PromQL

- Upitni jezik za Prometheus
- Dokumentacija: https://prometheus.io/docs/prometheus/latest/querying/basics/

## Primeri

### Broj obradjenih HTTP zahteva

    my_app_http_hit_total

### Broj odgovora sa status kodom 404

    response_status{status="404"}

### Broj GET zahteva po sekundi, za prethodnih 5 minuta

    rate(http_method{method="GET"}[5m])

### Prosecno vreme izvrsavanja zahteva za svaki endpoint

    http_response_time_sum / http_response_time_count