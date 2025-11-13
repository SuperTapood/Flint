from ..generated.test import K8STypes, K8SStack as K8SStack_, Stack, K8SConnection, ConnectionTypes, StackTypes
import sys
import socket

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
        socket_path = sys.argv[1]
    
        # Connect to Unix socket
        sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        sock.connect(socket_path)
        a = stack.SerializeToString()
        sock.sendall(a)
        sock.close()

