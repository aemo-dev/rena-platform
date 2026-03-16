/**
 * Loops blocks - repeat and iteration operations
 */
import * as Blockly from 'blockly';

// Loops Category
export const loopsCategory = {
  kind: 'category',
  name: 'Loops',
  colour: '120',
  contents: [
    {
      kind: 'block',
      type: 'controls_repeat_ext'
    },
    {
      kind: 'block',
      type: 'controls_whileUntil'
    },
    {
      kind: 'block',
      type: 'controls_for'
    },
    {
      kind: 'block',
      type: 'controls_forEach'
    },
    {
      kind: 'block',
      type: 'controls_flow_statements'
    }
  ]
};
