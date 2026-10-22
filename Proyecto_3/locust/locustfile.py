import random
from datetime import datetime, timezone
from locust import HttpUser, task, between

COUNTRIES = ["usa", "rus", "chn", "esp", "gtm"]

class MilitaryReportUser(HttpUser):
    wait_time = between(1, 3)
    endpoint = "/grpc-201700698"

    @task
    def send_report(self):
        payload = {
            "country": random.choice(COUNTRIES),
            "warplanes_in_air": random.randint(0, 50),
            "warships_in_water": random.randint(0, 30),
            "timestamp": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        }

        with self.client.post(self.endpoint, json=payload, catch_response=True) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"HTTP {response.status_code}: {response.text}")
