from ..generated import (
    K8STypes,
    K8SStack as K8SStack_,
    Stack,
    K8SConnection,
    ConnectionTypes,
    StackTypes,
)
import sys
import socket

# import traceback


# class AutoCollect:
#     def __init__(self):
#         self.items = []
#         self._prev_factory = None

#     def __enter__(self):
#         self._prev_factory = Element.factory
#         Element.factory = self.collecting_factory
#         return self.items

#     def __exit__(self, exc_type, exc, tb: traceback):
#         # Restore the original factory
#         Element.factory = self._prev_factory

#         if exc:
#             # Store the exception for later use
#             self.error = exc
#             print(exc)
#             print(exc_type)

#             # Decide if you want to suppress the exception:
#             # return True   → suppress
#             # return False  → re-raise
#             return True  # ← change to True if you want to swallow errors

#     def collecting_factory(self, *args, **kwargs):
#         obj = Element(*args, **kwargs)
#         self.items.append(obj)
#         return obj


# class Element:
#     def __init__(self, value):
#         self.value = value

#     @classmethod
#     def factory(cls, *a, **kw):
#         return cls(*a, **kw)


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
        k_stack = K8SStack_(objects=self.objects, namespace=self.namespace)
        k_conn = K8SConnection(api=self.api, token=self.token)
        stack = Stack(
            name=self.name,
            stack=StackTypes(k8s_stack=k_stack),
            connection=ConnectionTypes(k8s_connection=k_conn),
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
