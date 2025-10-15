from pyflint import *

pod = Pod(
    name="nginx",
    image="nginx:latest",
    ports=[
        80,
    ],
)
stack = K8SStack()
stack.add_objects(pod)

# with open("bob.bin", "wb") as file:
#     file.write(pod.SerializeToString())
stack.synth()
# print(len(open("bob.bin", "rb").read()))
