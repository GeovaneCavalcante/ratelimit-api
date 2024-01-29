from locust import HttpUser, task, between, constant

class RateLimitTest(HttpUser):
    # wait_time = constant(1)
    @task
    def check(self):
        headers = {'content-type': 'application/json'}
        headers = {'API_KEY': '5095bc00-2f9e-4e6f-b355-11688d20530d'}
        self.client.get("/health", headers=headers)