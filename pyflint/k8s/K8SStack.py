from ..generated import (
    K8STypes,
    K8SStack as K8SStack_,
    Stack,
    K8SConnection,
    ConnectionTypes,
    StackTypes,
    Secret,
    Pod,
    Deployment,
    Service
)
from ..common import BaseStack
import sys
from .K8SOutput import K8SOutput, K8STemplateOutput


class K8SStack(BaseStack):
    def __init__(self, api: str, token: str, name: str, namespace: str):
        """
        :param api: the api url for the kubernetes environment
        :param token: the token to use to authenticate against kubernetes
        """
        super().__init__()
        self.api = api
        self.token = token
        self.name = name
        self.namespace = namespace
        self.objects = []

    def add_objects(self, *objects):
        supported = [
            Secret,
            Deployment,
            Service,
            Pod,
            K8STemplateOutput,
        ]
        for obj in objects:
            if type(obj) not in supported:
                continue
            class_name = obj.__class__.__name__.lower()
            self.objects.append(K8STypes(**{class_name: obj}))

    def synth(self):
        k_stack = K8SStack_(objects=self.objects, namespace=self.namespace)
        k_conn = K8SConnection(api=self.api, token=self.token)
        stack = Stack(
            name=self.name,
            stack=StackTypes(k8sstack=k_stack),
            connection=ConnectionTypes(k8sconnection=k_conn),
        )
        self.send_data(stack.SerializeToString())

    def output(self, *args):
        if sys.version_info >= (3, 14):
            from string.templatelib import Template

            if type(args[0]) == Template:
                self.add_objects(K8STemplateOutput(args[0]))
                return
        self.add_objects(K8SOutput(*args))
