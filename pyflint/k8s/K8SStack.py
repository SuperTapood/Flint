from ..generated import (
    K8STypes,
    K8SStack as K8SStack_,
    Stack,
    K8SConnection,
    ConnectionTypes,
    StackTypes,
    Port,
    SecretData,
    ServiceTarget,
    K8SLookup,
)
from ..generated import k8s
import sys
import socket
from .K8SOutput import K8SOutput, K8STemplateOutput

from betterproto2 import Message
import inspect

import traceback


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
        self._prev_post_init = None

    def add_objects(self, *objects):
        excluded = [
            K8STypes,
            Port,
            SecretData,
            K8SStack,
            K8SStack_,
            K8SConnection,
            StackTypes,
            ConnectionTypes,
            Stack,
            ServiceTarget,
            K8SOutput,
            K8SLookup,
        ]
        for obj in objects:
            if type(obj) in excluded:
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
        if len(sys.argv) < 2:
            return
        socket_path = sys.argv[1]

        # Connect to Unix socket
        sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        sock.connect(socket_path)
        a = stack.SerializeToString()
        sock.sendall(a)
        sock.close()

    def output(self, *args):
        if sys.version_info >= (3, 14):
            from string.templatelib import Template

            if type(args[0]) == Template:
                self.add_objects(K8STemplateOutput(args[0]))
                return
        self.add_objects(K8SOutput(*args))

    def __enter__(self):
        self._prev_post_init = Message.__post_init__

        def post_init(obj):
            obj._unknown_fields = b""
            self.add_objects(obj)

        
        Message.__post_init__ = post_init

    def __exit__(self, exc_type, exc, tb: traceback):
        Message.__post_init__ = self._prev_post_init
        if exc:
            # Store the exception for later use
            self.error = exc
            print(exc)
            print(exc_type)

            # Decide if you want to suppress the exception:
            # return True   → suppress
            # return False  → re-raise
            return False  # ← change to True if you want to swallow errors
