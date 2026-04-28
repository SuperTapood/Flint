from os import listdir
import os
from pyflint import *
from subprocess import run


token = os.getenv("K8S_TOKEN")
# if token == None:
#     token = run(
#         ["kubectl", "create", "token", "flint"], capture_output=True
#     ).stdout.decode()


class Stack(K8SStack):
    def __init__(self):
        super().__init__(api="https://192.168.49.2:8443", token=token, name="wh", namespace="default")

        conatiner = Container(
            name="nginx",
            image="nginx:latest",
            ports=[
                80,
            ],
        )
        pod = Pod(
            containers=[
                conatiner,
            ],
        )
        
        self.daemon_set = DaemonSet(name="aset", replicas=1, pod=pod)

        port = Port(name="http", protocol="TCP", number=80)
        self.a = Deployment(name="a", replicas=1, pod=pod)
        # self.add_objects(deployment)
        self.deployment = Deployment(name="b", replicas=1, pod=pod)
        # self.add_objects(deployment)
        self.output(self.deployment["status"])
        service = Service(
            name="nginx-service",
            target=self.deployment,
            ports=[
                port,
            ],
        )
        # self.add_objects(service)
        self.secret = (Secret("my-secret", {"cock": "andballs"}))
        # self.output(
        #     t"My ip is {service['spec']['clusterIP']} {self.deployment['spec']['replicas']}"
        # )
        self.output("My replicas is ", "ADWw", self.deployment["spec"]["replicas"])
        self.ss = (
            StatefulSet(
                name="statefulnginx",
                replicas=1,
                pod=pod,
                volume_claim_templates=[
                    VolumeClaimTemplate(
                        name="www",
                        access_modes=[AccessMode.ReadWriteOnce],
                        storage_class_name="standard",
                        storage="1Gi",
                    )
                ],
            )
        )

        self.service = service


Stack().synth()
