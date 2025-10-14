from .generated.k8s.k8s_stack import _K8SStack as _stack
from .stack import Stack
import os


class K8SStack(Stack):
    def __init__(self):
        super().__init__()

        self.objects = []
    
    def add_objects(self, *objects):
        self.objects.extend(objects)
    
    def synth(self):
        stack = _stack(self.objects)
        print(stack.SerializeToString())
        print(stack.to_json())
        with open(f"{os.getcwd()}/bob.bin", "wb") as file:
            file.write(stack.SerializeToString())
