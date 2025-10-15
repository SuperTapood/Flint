from .generated.k8s.k8s_stack import _K8SStack, K8STypes
from .stack import Stack
from .generated.k8s.pod import Pod
import os


class K8SStack(Stack):
    def __init__(self):
        super().__init__()

        self.__objects = []
    
    def add_objects(self, *objects):
        for obj in objects:
            class_name = obj.__class__.__name__
            print(obj)
            self.__objects.append(K8STypes(**{class_name.lower(): obj}))
            print(self.__objects)
    
    def synth(self):
        stack = _K8SStack(self.__objects)
        print(stack.SerializeToString())
        print(stack.to_json())
        with open(f"{os.getcwd()}/bob.bin", "wb") as file:
            file.write(stack.SerializeToString())
