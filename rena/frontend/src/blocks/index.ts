/**
 * Main toolbox configuration for Blockly
 */
import { logicCategory } from './logic';
import { loopsCategory } from './loops';
import { mathCategory } from './math';
import { textCategory } from './text';
import { variablesCategory } from './variables';
import { functionsCategory } from './functions';
import { viewsCategory } from './views';

// Main toolbox definition
export const mainToolbox = {
  kind: 'categoryToolbox',
  contents: [
    viewsCategory,
    { kind: 'sep' },
    logicCategory,
    loopsCategory,
    mathCategory,
    textCategory,
    variablesCategory,
    functionsCategory
  ]
};

// Export all categories
export { logicCategory } from './logic';
export { loopsCategory } from './loops';
export { mathCategory } from './math';
export { textCategory } from './text';
export { variablesCategory } from './variables';
export { functionsCategory } from './functions';
export { viewsCategory } from './views';
export { ReactNativeGenerator } from './generators/react-native';
export { RenaTheme, applyBlockStyle, initializeTheme } from './renderers/rena-renderer';
