# /// script
# dependencies = [
#   "diagrams",
# ]
# ///

from diagrams import Cluster, Diagram, Edge
from diagrams.custom import Custom
from diagrams.k8s.network import Ingress, Service
from diagrams.onprem.client import Client
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.monitoring import Grafana, Prometheus

with Diagram():
    client = Client()

    with Cluster("Kubernetes Cluster"):
        service = Service()
        ingress = Ingress()

        client >> Edge(label="GET /r/{id}") >> ingress
        client >> Edge(label="POST /shorten") >> ingress
        ingress >> service

        metrics = Prometheus("metrics")
        metrics << Edge(color="firebrick", style="dashed") << Grafana("monitoring")

        shortener = Custom("url-shortener", "./docs/images/shortener.png")

        shortener >> PostgreSQL()
        service >> shortener << Edge(label="collect", color="firebrick", style="dashed") << metrics
