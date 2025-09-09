// Camel case conversion script
function toCamelCase(str) {
    str = String(str);
    return str.replace(/_([a-z])/g, (g) => g[1].toUpperCase())
              .replace(/^[A-Z]/, (g) => g.toLowerCase());
}

// input is already parsed JSON object
let result = "export interface " + input.tableName + " {\n";

// Process each field
for (const field of input.fields) {
    const fieldName = String(field.name || '');
    const camelCaseName = toCamelCase(fieldName);

    if (field.comment) {
        result += "  /** " + field.comment + " */\n";
    }
    result += "  " + camelCaseName + ": " + field.tsType + ";\n";
}

result += "}\n";

// Set output result
output = result;