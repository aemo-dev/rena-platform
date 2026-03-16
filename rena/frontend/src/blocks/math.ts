/**
 * Math blocks - mathematical operations
 */
import * as Blockly from 'blockly';

// Math Category
export const mathCategory = {
  kind: 'category',
  name: 'Math',
  colour: '290',
  contents: [
    {
      kind: 'block',
      type: 'math_number'
    },
    {
      kind: 'block',
      type: 'math_arithmetic'
    },
    {
      kind: 'block',
      type: 'math_single'
    },
    {
      kind: 'block',
      type: 'math_trig'
    },
    {
      kind: 'block',
      type: 'math_constant'
    },
    {
      kind: 'block',
      type: 'math_number_property'
    },
    {
      kind: 'block',
      type: 'math_round'
    },
    {
      kind: 'block',
      type: 'math_on_list'
    },
    {
      kind: 'block',
      type: 'math_modulo'
    },
    {
      kind: 'block',
      type: 'math_constrain'
    },
    {
      kind: 'block',
      type: 'math_random_int'
    },
    {
      kind: 'block',
      type: 'math_random_float'
    }
  ]
};
