from ..generated import K8SOutput as _output, K8SLookup, K8STypes
import sys
import uuid

from betterproto2 import Message

Message.getitem = lambda self, item: K8SLookup(
    object=K8STypes(**{self.__class__.__name__.lower(): self}),
    keys=[
        item,
    ],
)

K8SLookup.getitem = lambda self, item: [self.keys.append(item), self][1]

if sys.version_info >= (3, 14):
    from string.templatelib import Template

    def K8STemplateOutput(template: Template):
        lookups = []
        for inter in template.interpolations:
            if type(inter.value) == K8SLookup:
                lookups.append(inter.value)
            else:
                lookups.append(
                    K8SLookup(
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

    def K8STemplateOutput(template: any):
        raise NotImplementedError(
            "Output requires Python 3.14 or higher, use OldOutput if you run older versions"
        )


def K8SOutput(*values):
    lookups = []
    strings = []
    for val in values:
        if type(val) == str:
            strings.append(val)
        else:
            lookups.append(val)

    o = _output(lookups=lookups, strings=strings, id=uuid.uuid1().__str__())
    return o
