/**
 * React Native code generator
 */
import * as Blockly from 'blockly';

// Define the React Native generator
export const ReactNativeGenerator = new Blockly.Generator('ReactNative');

// Order of operations
(ReactNativeGenerator as any).ORDER_ATOMIC = 0;
(ReactNativeGenerator as any).ORDER_NONE = 0;

// Initialize the generator
ReactNativeGenerator.init = function(workspace: Blockly.Workspace) {
  (this as any).imports = new Set<string>();
  (this as any).componentImports = new Map<string, Set<string>>();
  (this as any).variables = Blockly.Variables.allUsedVarModels(workspace);
};

// Finish the generator
ReactNativeGenerator.finish = function(code: string) {
  let importCode = "import React, { useState, useEffect } from 'react';\n";
  
  // Group imports by package
  (this as any).componentImports.forEach((tags: Set<string>, pkg: string) => {
    const tagList = Array.from(tags).join(', ');
    importCode += `import { ${tagList} } from '${pkg}';\n`;
  });

  let varCode = "";
  (this as any).variables.forEach((v: any) => {
    varCode += `  const [${v.name}, set${v.name.charAt(0).toUpperCase() + v.name.slice(1)}] = useState(null);\n`;
  });

  return `${importCode}\n\nexport default function App() {\n${varCode}\n  return (\n${code}\n  );\n}`;
};

// Generator for App Main View block
(ReactNativeGenerator as any)['app_main_view'] = function(block: Blockly.Block) {
  const viewName = block.getFieldValue('VIEW_NAME') || 'Home';
  const children = ReactNativeGenerator.statementToCode(block, 'CONTENT');
  
  return `    <View style={{flex: 1}}>\n      {/* ${viewName} View */}
${children}    </View>\n`;
};

// Basic Logic Blocks
(ReactNativeGenerator as any)['controls_if'] = function(block: Blockly.Block) {
  const argument0 = ReactNativeGenerator.valueToCode(block, 'IF0', (ReactNativeGenerator as any).ORDER_NONE) || 'false';
  const branch0 = ReactNativeGenerator.statementToCode(block, 'DO0');
  let code = `{${argument0} && (\n${branch0}\n)}`;
  return code;
};
