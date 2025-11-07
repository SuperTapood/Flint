from .generated.k8s.k8s_stack_ import K8STypes, K8S_Stack_
from .generated.common.stack_ import Stack, StackTypes, ConnectionTypes
from .generated.k8s.k8s_connection import K8S_Connection
import sys
import os


class K8SStack:
    def __init__(self, api: str, token: str, name: str, namespace: str):
        """
        :param api: the api url for the kubernetes environment
        :param token: the token to use to authenticate against kubernetes
        """
        self.api = api
        self.token = token
        self.name = name
        self.namespace = namespace
        self.objects = []

    def add_objects(self, *objects):
        for obj in objects:
            class_name = obj.__class__.__name__.lower()
            self.objects.append(K8STypes(**{class_name: obj}))

    def synth(self):
        k_stack = K8S_Stack_(objects=self.objects, namespace=self.namespace)
        k_conn = K8S_Connection(api=self.api, token=self.token)
        stack = Stack(
            name=self.name,
            stack=StackTypes(k8s_stack=k_stack),
            connection=ConnectionTypes(k8s_connection=k_conn),
        )
        if len(sys.argv) > 1 and sys.argv[1].isdigit():
            fd = int(sys.argv[1])
            with os.fdopen(fd, "wb") as file:
                file.write(stack.SerializeToString())
                file.flush()
