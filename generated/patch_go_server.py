import json

open_api_schema = "api/openapi.json"

with open(open_api_schema) as f:
    schema = json.load(f)

authorizePaths = set()

for path, path_val in schema["paths"].items():
    for method, method_val in path_val.items():
        if "security" in method_val:
            authorizePaths.add((path, method))


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
