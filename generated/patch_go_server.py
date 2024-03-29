import json
import sys

with open(sys.argv[1]) as f:
    schema = json.load(f)

authorizePaths = list()

for path, path_val in schema["paths"].items():
    for method, method_val in path_val.items():
        if "security" in method_val:
            authorizePaths.append((path, method))


print("""package openapi

var AuthorizeRoutes = []struct{
    Path   string
    Method string
} {""")
for paths in authorizePaths:
    print("""
    {{
        Path:   "{}",
        Method: "{}",
    }},""".format(paths[0], paths[1].upper()))

print("""
}

""")
