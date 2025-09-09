// Code formatting script
// input is already parsed JSON object

// Get current generated code
let result = tsCode;

// Format code - add proper indentation and line breaks
result = result
    .replace(/\{\s*\n/g, '{\n  ')
    .replace(/;\s*\n/g, ';\n  ')
    .replace(/\n\s*\}/g, '\n}')
    .replace(/\n\s*\n/g, '\n')
    .trim();

// Ensure there's a newline at the end
if (!result.endsWith('\n')) {
    result += '\n';
}

// Set output result
output = result;