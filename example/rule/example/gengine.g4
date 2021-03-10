grammar gengine;

primary: ruleEntity+;

ruleEntity: RULE ruleName ruleDesc?salience? BEGIN ruleContent END;
ruleName : testStr;
ruleDesc : testStr;

testStr: ' ';
salience : ;
ruleContent: 'return 1';