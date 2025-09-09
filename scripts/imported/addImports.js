// Type import script
// input is already parsed JSON object

// Get current generated code
let result = tsCode;

// Collect types that need to be imported
const importTypes = new Set();

// Check field types to determine what needs to be imported
for (const field of input.fields) {
    const type = field.tsType.toLowerCase();
    if (type.includes('date') || type.includes('time')) {
        importTypes.add('Date');
    }
}

// If there are types to import, add import statements
if (importTypes.size > 0) {
    let imports = '// Auto-generated import statements\n';
    for (const type of importTypes) {
        imports += `// import { ${type} } from './types'\n`;
    }
    imports += '\n';
    
    // Add import statements at the beginning of the code
    result = imports + result;
}

// Set output result
output = result;