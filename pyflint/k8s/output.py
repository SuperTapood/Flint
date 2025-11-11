from ..generated.k8s.k8s_stack_ import K8SOutput as _output
import sys

if sys.version_info >= (3, 14):
    class Lookup:
        def __init__(self, obj, index=None):
            self.object = obj
            self.indices = [index, ] if index is not None else [] 

        def __getitem__(self, index):
            self.indices.append(index)
            return self
        
    from betterproto import Message

    Message.__getitem__ = lambda self, item: Lookup(self, item)
        
    from string.templatelib import Template
    def Output(template: Template):
        objects = []
        indices = []
        strings = []
        for string, inter in zip(template.strings, template.interpolations):
            strings.append(string)
            objects.append(inter.value.object)
            indices.append(inter.value.indices)
        return _output(objects=objects, indices=indices, strings=strings)
else:
    def Output(template: any):
        raise NotImplementedError("Output requires Python 3.14 or higher, use OldOutput if you run older versions")