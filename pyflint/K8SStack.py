from .generated.k8s.k8s_stack_ import K8STypes, K8S_Stack_
from .generated.common.stack_ import Stack_, StackTypes

class K8SStack:
    def __init__(self):
        self.objects = []
    
    def add_objects(self, *objects):
        for obj in objects:
            class_name = obj.__class__.__name__.lower()
            self.objects.append(K8STypes(**{class_name: obj}))
    
    def synth(self):
        k_stack = K8S_Stack_(self.objects)
        stack = Stack_(StackTypes(k8s_stack=k_stack))
        print(stack)
        with open("bob.bin", "wb") as file:
            file.write(stack.SerializeToString())