from pyflint.generated.thing import Person

a = Person(name="bob", id=1234, email="bob@example.com")

with open("bob.bin", "wb") as f:
    f.write(a.SerializeToString())