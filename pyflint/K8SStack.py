from .generated.k8s.k8s_stack_ import K8STypes, K8S_Stack_
from .generated.common.stack_ import Stack_, StackTypes
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
        k_stack = K8S_Stack_(
            self.objects, self.api, self.token, self.name, self.namespace
        )
        stack = Stack_(StackTypes(k8s_stack=k_stack))
        if len(sys.argv) > 1 and sys.argv[1].isdigit():
            fd = int(sys.argv[1])
            with os.fdopen(fd, "wb") as file:
                file.write(stack.SerializeToString())
                file.flush()
            # print(str(stack.SerializeToString())[1:-2])
            # with os.fdopen(1, "wb", closefd=False) as stdout:
            #     stdout.write(stack.SerializeToString())
            #     stdout.flush()
        # with open("flintcore/bib.bin", "wb") as file:
        #     file.write(stack.SerializeToString())
