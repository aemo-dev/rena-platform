/**
 * Custom Blockly Theme - Rena Platform Style
 * Inspired by MIT App Inventor's block styling
 */
import * as Blockly from 'blockly/core';

/**
 * Define the Rena Platform theme
 */
export const RenaTheme: Blockly.Theme = Blockly.Theme.defineTheme('rena', {
  'name': 'Rena',
  'base': Blockly.Themes.Classic,
  'blockStyles': {
    // Views category - Purple/Pink
    'view_blocks': {
      'colourPrimary': '#9C27B0',
      'colourSecondary': '#7B1FA2',
      'colourTertiary': '#E1BEE7',
      'hat': 'cap'
    },
    // Components category - Blue
    'component_blocks': {
      'colourPrimary': '#2196F3',
      'colourSecondary': '#1976D2',
      'colourTertiary': '#BBDEFB'
    },
    // Logic category - Green
    'logic_blocks': {
      'colourPrimary': '#4CAF50',
      'colourSecondary': '#388E3C',
      'colourTertiary': '#C8E6C9'
    },
    // Loops category - Orange
    'loop_blocks': {
      'colourPrimary': '#FF9800',
      'colourSecondary': '#F57C00',
      'colourTertiary': '#FFE0B2'
    },
    // Math category - Red
    'math_blocks': {
      'colourPrimary': '#F44336',
      'colourSecondary': '#D32F2F',
      'colourTertiary': '#FFCDD2'
    },
    // Text category - Softer Gold (not too bright)
    'text_blocks': {
      'colourPrimary': '#F9A825',  // Darker gold
      'colourSecondary': '#F57F17',
      'colourTertiary': '#FFECB3'
    },
    // Variables category - Cyan
    'variable_blocks': {
      'colourPrimary': '#00BCD4',
      'colourSecondary': '#0097A7',
      'colourTertiary': '#B2EBF2'
    },
    // Functions category - Indigo
    'procedure_blocks': {
      'colourPrimary': '#3F51B5',
      'colourSecondary': '#303F9F',
      'colourTertiary': '#C5CAE9'
    }
  },
  'categoryStyles': {
    'view_category': {
      'colour': '#AB47BC'  // Softer purple
    },
    'component_category': {
      'colour': '#42A5F5'  // Softer blue
    },
    'logic_category': {
      'colour': '#66BB6A'  // Softer green
    },
    'loop_category': {
      'colour': '#FFA726'  // Softer orange
    },
    'math_category': {
      'colour': '#EF5350'  // Softer red
    },
    'text_category': {
      'colour': '#F9A825'  // Softer gold (matches block)
    },
    'variable_category': {
      'colour': '#26C6DA'  // Softer cyan
    },
    'procedure_category': {
      'colour': '#5C6BC0'  // Softer indigo
    }
  },
  'fontStyle': {
    'family': '"Inter Variable", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    'size': 12,
    'weight': '500'
  }
});

/**
 * Helper function to apply custom styles to blocks based on their type
 */
export function applyBlockStyle(block: Blockly.Block) {
  const type = block.type;
  
  // Determine style based on block type prefix or category
  if (type.startsWith('app_') || type.includes('view')) {
    block.setStyle('view_blocks');
  } else if (type.startsWith('component_')) {
    block.setStyle('component_blocks');
  } else if (type.startsWith('controls_')) {
    block.setStyle('loop_blocks');
  } else if (type.startsWith('logic_')) {
    block.setStyle('logic_blocks');
  } else if (type.startsWith('math_')) {
    block.setStyle('math_blocks');
  } else if (type.startsWith('text_')) {
    block.setStyle('text_blocks');
  } else if (type === 'variables_get' || type === 'variables_set') {
    block.setStyle('variable_blocks');
  } else if (type.startsWith('procedures_')) {
    block.setStyle('procedure_blocks');
  }
}

/**
 * Initialize the theme with workspace
 */
export function initializeTheme(workspace: Blockly.WorkspaceSvg) {
  workspace.setTheme(RenaTheme);
  
  // Apply styles to existing blocks
  workspace.getAllBlocks(false).forEach(block => {
    applyBlockStyle(block);
  });
}
