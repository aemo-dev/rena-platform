/**
 * Main views - app structure blocks
 */
import * as Blockly from 'blockly';

// Views Category
export const viewsCategory = {
  kind: 'category',
  name: 'Views',
  colour: '290',
  contents: [
    {
      kind: 'block',
      type: 'app_main_view',
      enabled: true
    }
  ]
};

// Register App Main View block
Blockly.Blocks['app_main_view'] = {
  init: function(this: any) {
    this.setNextStatement(true, null);
    
    this.appendDummyInput()
      .appendField("Main View")
      .appendField(new Blockly.FieldTextInput("Home"), "VIEW_NAME");
    this.appendStatementInput("CONTENT")
      .setCheck(null)
      .appendField("Content");
    this.setColour(290);
    this.setTooltip("Main view/screen of your React Native app");
    this.setHelpUrl("");
  }
};
