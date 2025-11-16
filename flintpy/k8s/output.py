from ..generated import K8SOutput as _output, Lookup, K8STypes
import sys
import uuid

if sys.version_info >= (3, 14):
    from betterproto2 import Message

    Message.getitem = lambda self, item: Lookup(
        object=K8STypes(**{self.__class__.__name__.lower(): self}),
        keys=[
            item,
        ],
    )

    Lookup.getitem = lambda self, item: [self.keys.append(item), self][1]

    from string.templatelib import Template

    def Output(template: Template):
        lookups = []
        for inter in template.interpolations:
            if type(inter.value) == Lookup:
                lookups.append(inter.value)
            else:
                lookups.append(
                    Lookup(
                        object=K8STypes(
                            **{inter.value.__class__.__name__.lower(): inter.value}
                        ),
                        keys=[],
                    )
                )
        strings = [string for string in template.strings]
        o = _output(lookups=lookups, strings=strings, id=uuid.uuid8().__str__())
        return o
else:

    def Output(template: any):
        raise NotImplementedError(
            "Output requires Python 3.14 or higher, use OldOutput if you run older versions"
        )
