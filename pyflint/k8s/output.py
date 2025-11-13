from ..generated.test import K8SOutput as _output, Lookup, K8STypes
import sys

if sys.version_info >= (3, 14):
    from betterproto2 import Message

    Message.getitem = lambda self, item: Lookup(K8STypes(**{self.__class__.__name__.lower(): self}), item)
        
    from string.templatelib import Template
    def Output(template: Template):
        lookups = []
        strings = []
        for string, inter in zip(template.strings, template.interpolations):
            strings.append(string)
            lookups.append(inter.value.object)
        return _output(lookups=lookups, strings=strings)
else:
    def Output(template: any):
        raise NotImplementedError("Output requires Python 3.14 or higher, use OldOutput if you run older versions")